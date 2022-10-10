package mongodb

import (
	"bytes"
	"errors"
	"strings"

	"scraper/src/types"

	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"

	"time"

	"path/filepath"

	"image"
	"image/jpeg"
	"image/png"
)

// InsertImage insert an image in its collection
func InsertImage(collection *mongo.Collection, image types.Image) (primitive.ObjectID, error) {
	res, err := collection.InsertOne(context.TODO(), image)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("safecast of ObjectID did not work")
	}
	return insertedID, nil
}

// RemoveImage remove an image based on its mongodb id
func RemoveImage(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// RemoveImageAndFile remove an image based on its mongodb id and remove its file
func RemoveImageAndFile(s3Client *s3.Client, collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	image, err := FindOne[types.Image](collection, bson.M{"_id": id})
	if err != nil {
		return nil, fmt.Errorf("FindImageByID has failed: %v", err)
	}
	deletedCount, err := RemoveImage(collection, id)
	if err != nil {
		return nil, fmt.Errorf("RemoveImage has failed: %v", err)
	}

	path := filepath.Join(image.Origin, image.Name)

	_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(utils.GetEnvVariable("IMAGES_BUCKET")),
		Key:    aws.String(path),
	})

	// sometimes images can have the same file stored but are present multiple in the search request
	if err != nil && *deletedCount == 0 {
		return nil, fmt.Errorf("os.Remove has failed: %v", err)
	}
	return deletedCount, nil
}

func RemoveImagesAndFiles(s3Client *s3.Client, mongoClient *mongo.Client, query bson.M, options *options.FindOptions) (*int64, error) {
	collectionImages := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PRODUCTION"))
	var deletedCount int64
	images, err := FindMany[types.Image](collectionImages, query, options)
	if err != nil {
		return nil, fmt.Errorf("FindImagesIDs has failed: %v", err)
	}
	for _, image := range images {
		deletedOne, err := RemoveImageAndFile(s3Client, collectionImages, image.ID)
		if err != nil {
			return nil, fmt.Errorf("RemoveImageAndFile has failed for %s: %v", image.ID.Hex(), err)
		}
		deletedCount += *deletedOne
	}
	return &deletedCount, nil
}

// UpdateImageTags add tags to an image based on its mongodb id
func UpdateImageTagsPush(collection *mongo.Collection, body types.BodyUpdateImageTagsPush) (*int64, error) {
	query := bson.M{"_id": body.ID}
	for i := 0; i < len(body.Tags); i++ {
		tag := &body.Tags[i]
		now := time.Now()
		tag.CreationDate = &now
	}
	update := bson.M{"$push": bson.M{"tags": bson.M{"$each": body.Tags}}}
	res, err := collection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return nil, fmt.Errorf("UpdateOne has failed: %v", err)
	}
	return &res.ModifiedCount, nil
}

// UpdateImageTagsPull removes specific tags from an image
func UpdateImageTagsPull(collection *mongo.Collection, body types.BodyUpdateImageTagsPull) (*int64, error) {
	query := bson.M{
		"_id":    body.ID,
		"origin": body.Origin,
	}
	update := bson.M{
		"$pull": bson.M{
			"tags": bson.M{
				"name": bson.M{
					"$in": body.Names,
				},
			},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return nil, err
	}
	return &res.ModifiedCount, nil
}

// cropFileAndData updates the data in db and crop the original file
func cropFileAndData(s3Client *s3.Client, mongoCollection *mongo.Collection, body types.BodyImageCrop) (image.Image, *types.Image, error) {
	// get information of the image
	query := bson.M{"_id": body.ID}
	options := options.FindOne().SetProjection(bson.M{"_id": 0}) // no _id for replacing later on
	imageData, err := FindOne[types.Image](mongoCollection, query, options)
	if err != nil {
		return nil, nil, fmt.Errorf("FindOne[Image] has failed: %v", err)
	}

	// update the image size and tags boxes
	imageData, err = updateImageBoxes(body, imageData)
	if err != nil {
		return nil, nil, fmt.Errorf("updateImageBoxes has failed: %v", err)
	}

	path := filepath.Join(imageData.Origin, imageData.Name)
	buffer, err := utils.GetItemS3(s3Client, path)
	if err != nil {
		return nil, nil, err
	}

	// convert []byte to image
	img, _, _ := image.Decode(bytes.NewReader(buffer))

	// crop the image with the bounding box rectangle
	cropRect := image.Rect(*body.Box.Tlx, *body.Box.Tly, *body.Box.Tlx+*body.Box.Width, *body.Box.Tly+*body.Box.Height)
	img, err = utils.CropImage(img, cropRect)
	if err != nil {
		return nil, nil, err
	}

	return img, imageData, nil
}

// UpdateImageFile update the image with its tags when it is cropped
func UpdateImageCrop(s3Client *s3.Client, mongoClient *mongo.Client, body types.BodyImageCrop) (*int64, error) {
	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))

	// crop data and file
	img, imageData, err := cropFileAndData(s3Client, collectionImagesPending, body)
	if err != nil {
		return nil, fmt.Errorf("cropFileAndData has failed: %v", err)
	}

	// replace in db and file of the updated image
	updatedCount, err := replaceImage(s3Client, collectionImagesPending, imageData, img)
	if err != nil {
		return nil, fmt.Errorf("replaceImage has failed: %v", err)
	}
	return updatedCount, nil
}

func CreateImageCrop(s3Client *s3.Client, mongoClient *mongo.Client, body types.BodyImageCrop) (*int64, error) {
	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))

	// crop data and file
	img, imageData, err := cropFileAndData(s3Client, collectionImagesPending, body)
	if err != nil {
		return nil, fmt.Errorf("cropFileAndData has failed: %v", err)
	}

	// add the current date and time to the new name
	imageData.Name = fmt.Sprintf("%s_%s", imageData.OriginID, time.Now().Format(time.RFC3339))

	// replace in db and file of the updated image
	updatedCount, err := replaceImage(s3Client, collectionImagesPending, imageData, img)
	if err != nil {
		return nil, fmt.Errorf("replaceImage has failed: %v", err)
	}
	return updatedCount, nil
}

func CopyImage(s3Client *s3.Client, mongoClient *mongo.Client, body types.BodyImageCopy) (*string, error) {
	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))

	sourcePath := fmt.Sprintf("%s/%s.%s", body.Origin, body.Name, body.Extension)
	name := fmt.Sprintf("%s_%s", body.OriginID, time.Now().Format(time.RFC3339))
	destinationPath := fmt.Sprintf("%s/%s.%s", body.Origin, name, body.Extension)

	query := bson.M{"name": body.Name}
	image, err := FindOne[types.Image](collectionImagesPending, query)
	if err != nil {
		return nil, fmt.Errorf("FindOne[Image] has failed: %v", err)
	}
	image.ID = primitive.NilObjectID
	image.Name = name
	_, err = collectionImagesPending.InsertOne(context.TODO(), *image)
	if err != nil {
		return nil, fmt.Errorf("InsertOne has failed: %v", err)
	}
	return utils.CopyItemS3(s3Client, sourcePath, destinationPath)
}

func updateImageBoxes(body types.BodyImageCrop, imageData *types.Image) (*types.Image, error) {
	// new size creation
	imageSizeID := primitive.NewObjectID()
	now := time.Now()
	size := types.ImageSize{
		ID:           imageSizeID,
		CreationDate: &now,
		Box:          body.Box, // absolute position
	}
	imageData.Size = append(imageData.Size, size)

	i := 0
	for {
		if i >= len(imageData.Tags) {
			break
		}
		tag := imageData.Tags[i]
		if (types.Box{}) != tag.Origin.Box {
			// relative position of tags
			tlx := *tag.Origin.Box.Tlx
			tly := *tag.Origin.Box.Tly
			width := *tag.Origin.Box.Width
			height := *tag.Origin.Box.Height

			// box outside on the image right
			if tlx > *body.Box.Tlx+*body.Box.Width {
				// last element removed
				if i == len(imageData.Tags)-1 {
					imageData.Tags = imageData.Tags[:i]
				} else { // not last element removed
					imageData.Tags = append(imageData.Tags[:i], imageData.Tags[i+1:]...)
				}
				continue
			}
			// box left outside on the image left
			if tlx < *body.Box.Tlx {
				// box outside on the image left
				if tlx+width < *body.Box.Tlx {
					width = 0
				} else { // box right inside the image
					width = width - *body.Box.Tlx + tlx
				}
				tlx = *body.Box.Tlx
			} else { // box left inside image
				// box right outside on the image right
				if tlx+width > *body.Box.Tlx+*body.Box.Width {
					width = *body.Box.Tlx + *body.Box.Width - tlx
				}
				tlx = tlx - *body.Box.Tlx
			}
			// box width too small
			if width < 50 {
				// last element removed
				if i == len(imageData.Tags)-1 {
					imageData.Tags = imageData.Tags[:i]
				} else { // not last element removed
					imageData.Tags = append(imageData.Tags[:i], imageData.Tags[i+1:]...)
				}
				continue
			}

			// box outside at the image bottom
			if tly > *body.Box.Tly+*body.Box.Height {
				// last element removed
				if i == len(imageData.Tags)-1 {
					imageData.Tags = imageData.Tags[:i]
				} else { // not last element removed
					imageData.Tags = append(imageData.Tags[:i], imageData.Tags[i+1:]...)
				}
				continue
			}
			// box top outside on the image top
			if tly < *body.Box.Tly {
				// box outside on the image top
				if tly+height < *body.Box.Tly {
					height = 0
				} else { // box bottom inside the image
					height = height - *body.Box.Tly + tly
				}
				tly = *body.Box.Tly
			} else { // box top inside image
				// box bottom outside on the image bottom
				if tly+height > *body.Box.Tly+*body.Box.Height {
					height = *body.Box.Tly + *body.Box.Height - tly
				}
				tly = tly - *body.Box.Tly
			}
			// box height too small
			if height < 50 {
				// last element removed
				if i == len(imageData.Tags)-1 {
					imageData.Tags = imageData.Tags[:i]
				} else { // not last element removed
					imageData.Tags = append(imageData.Tags[:i], imageData.Tags[i+1:]...)
				}
				continue
			}

			// set the new relative reference to the newly cropped image
			tag.Origin.ImageSizeID = imageSizeID
			tag.Origin.Box.Tlx = &tlx
			tag.Origin.Box.Tly = &tly
			tag.Origin.Box.Width = &width
			tag.Origin.Box.Height = &height
		}
		i++
	}
	return imageData, nil
}

func replaceImage(s3Client *s3.Client, collection *mongo.Collection, imageData *types.Image, img image.Image) (*int64, error) {
	// update in db the new dimensions, tag boxes and new size
	query := bson.M{"name": imageData.Name} // match the existing or new name
	options := options.Replace().SetUpsert(true)
	res, err := collection.ReplaceOne(context.TODO(), query, imageData, options)
	if err != nil {
		return nil, fmt.Errorf("UpdateOne has failed: %v", err)
	}
	if res.UpsertedCount == 0 && res.ModifiedCount == 0 {
		return nil, fmt.Errorf("no upsert or update have been done")
	}

	// create buffer
	buffer := new(bytes.Buffer)
	// encode image to buffer

	if imageData.Extension == "jpeg" || imageData.Extension == "jpg" {
		err := jpeg.Encode(buffer, img, nil)
		if err != nil {
			return nil, fmt.Errorf("jpeg.Encode has failed: %v", err)
		}
	} else if imageData.Extension == "png" {
		err := png.Encode(buffer, img)
		if err != nil {
			return nil, fmt.Errorf("png.Encode has failed: %v", err)
		}
	} else {
		return nil, fmt.Errorf("no image extension matching the buffer conversion")
	}

	// convert buffer to reader
	reader := bytes.NewReader(buffer.Bytes())

	// upload new image in s3
	path := fmt.Sprintf("%s/%s.%s", imageData.Origin, imageData.Name, imageData.Extension)
	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(utils.GetEnvVariable("IMAGES_BUCKET")),
		Key:    aws.String(path),
		Body:   reader,
	})
	if err != nil {
		return nil, fmt.Errorf("uploader.Upload has failed: %v", err)
	}

	return &res.ModifiedCount, nil
}

// InsertImageUnwanted insert an unwanted image
func InsertImageUnwanted(mongoClient *mongo.Client, body types.Image) (interface{}, error) {
	now := time.Now()
	body.CreationDate = &now
	body.Origin = strings.ToLower(body.Origin)

	// insert the unwanted image
	collectionImagesUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("UNDESIRED"))
	res, err := collectionImagesUnwanted.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, fmt.Errorf("InsertOne has failed: %v", err)
	}
	return res.InsertedID, nil
}

func TransferImage(mongoClient *mongo.Client, body types.BodyTransferImage) (interface{}, error) {
	collectionImagesFrom, err := utils.ImagesCollection(mongoClient, body.From)
	if err != nil {
		return nil, err
	}
	collectionImagesTo, err := utils.ImagesCollection(mongoClient, body.To)
	if err != nil {
		return nil, err
	}
	query := bson.M{"originID": body.OriginID}
	image, err := FindOne[types.Image](collectionImagesFrom, query)
	if err != nil {
		return nil, fmt.Errorf("FindOne[Image] has failed: %v", err)
	}
	image.ID = primitive.NilObjectID
	res, err := collectionImagesTo.InsertOne(context.TODO(), *image)
	if err != nil {
		return nil, fmt.Errorf("InsertOne has failed: %v", err)
	}
	_, err = collectionImagesFrom.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, fmt.Errorf("DeleteOne has failed: %v", err)
	}
	return res.InsertedID, nil
}
