package mongodb

import (
	"errors"
	"strings"

	"scraper/src/types"

	"scraper/src/utils"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"

	"time"

	"path/filepath"

	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// InsertImage insert an image in its collection
func InsertImage(collection *mongo.Collection, image types.Image) (primitive.ObjectID, error) {
	res, err := collection.InsertOne(context.TODO(), image)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("Safecast of ObjectID did not work")
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
func RemoveImageAndFile(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	image, err := FindOne[types.Image](collection, bson.M{"_id": id})
	if err != nil {
		return nil, fmt.Errorf("FindImageByID has failed: %v", err)
	}
	deletedCount, err := RemoveImage(collection, id)
	if err != nil {
		return nil, fmt.Errorf("RemoveImage has failed: %v", err)
	}
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := filepath.Join(folderDir, image.Origin, image.Name)
	err = os.Remove(path)
	// sometimes images can have the same file stored but are present multiple in the search request
	if err != nil && *deletedCount == 0 {
		return nil, fmt.Errorf("os.Remove has failed: %v", err)
	}
	return deletedCount, nil
}

func RemoveImagesAndFiles(mongoClient *mongo.Client, query bson.M, options *options.FindOptions) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_WANTED_COLLECTION"))
	var deletedCount int64
	images, err := FindMany[types.Image](collectionImages, query, options)
	if err != nil {
		return nil, fmt.Errorf("FindImagesIDs has failed: %v", err)
	}
	for _, image := range images {
		deletedOne, err := RemoveImageAndFile(collectionImages, image.ID)
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
	update := bson.M{
		"$push": bson.M{
			"tags": bson.M{"$each": body.Tags},
		},
	}
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

// UpdateImageFile update the image with its tags when it is cropped
func UpdateImageCrop(mongoClient *mongo.Client, body types.BodyImageCrop) (*int64, error) {
	collectionImagesPending := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_PENDING_COLLECTION"))
	// update the image size and tags boxes
	image, err := updateImageBoxes(collectionImagesPending, body)
	if err != nil {
		return nil, fmt.Errorf("updateImageBoxes has failed: %v", err)
	}
	// replace in db and file of the updated image
	updatedCount, err := replaceImage(collectionImagesPending, image, body.File)
	if err != nil {
		return nil, fmt.Errorf("replaceImage has failed: %v", err)
	}
	return updatedCount, nil
}

// UpdateImageFile update the image with its tags when it is cropped
func CreateImageCrop(mongoClient *mongo.Client, body types.BodyImageCrop) (*int64, error) {
	collectionImagesPending := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_PENDING_COLLECTION"))
	// update the image size and tags boxes
	image, err := updateImageBoxes(collectionImagesPending, body)
	if err != nil {
		return nil, fmt.Errorf("updateImageBoxes has failed: %v", err)
	}
	// add the current date and time to the name
	image.Name = fmt.Sprintf("%s_%s.%s", image.ID, time.Now().Format(time.RFC3339), image.Extension)
	// replace in db and file of the updated image
	updatedCount, err := replaceImage(collectionImagesPending, image, body.File)
	if err != nil {
		return nil, fmt.Errorf("replaceImage has failed: %v", err)
	}
	return updatedCount, nil
}

func updateImageBoxes(collection *mongo.Collection, body types.BodyImageCrop) (*types.Image, error) {
	query := bson.M{"_id": body.ID}
	imageFound, err := FindOne[types.Image](collection, query)
	if err != nil {
		return nil, fmt.Errorf("FindOne[Image] has failed: %v", err)
	}

	// new size creation
	imageSizeID := primitive.NewObjectID()
	now := time.Now()
	size := types.ImageSize{
		ID:           imageSizeID,
		CreationDate: &now,
		Box:          body.Box, // absolute position
	}
	imageFound.Size = append(imageFound.Size, size)

	i := 0
	for {
		if i >= len(imageFound.Tags) {
			break
		}
		tag := imageFound.Tags[i]
		if (types.Box{}) != tag.Origin.Box {
			// relative position of tags
			tlx := *tag.Origin.Box.X
			tly := *tag.Origin.Box.Y
			width := *tag.Origin.Box.Width
			height := *tag.Origin.Box.Height

			// box outside on the image right
			if tlx > *body.Box.X+*body.Box.Width {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}
			// box left outside on the image left
			if tlx < *body.Box.X {
				// box outside on the image left
				if tlx+width < *body.Box.X {
					width = 0
				} else { // box right inside the image
					width = width - *body.Box.X + tlx
				}
				tlx = *body.Box.X
			} else { // box left inside image
				// box right outside on the image right
				if tlx+width > *body.Box.X+*body.Box.Width {
					width = *body.Box.X + *body.Box.Width - tlx
				}
				tlx = tlx - *body.Box.X
			}
			// box width too small
			if width < 50 {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}

			// box outside at the image bottom
			if tly > *body.Box.Y+*body.Box.Height {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}
			// box top outside on the image top
			if tly < *body.Box.Y {
				// box outside on the image top
				if tly+height < *body.Box.Y {
					height = 0
				} else { // box bottom inside the image
					height = height - *body.Box.Y + tly
				}
				tly = *body.Box.Y
			} else { // box top inside image
				// box bottom outside on the image bottom
				if tly+height > *body.Box.Y+*body.Box.Height {
					height = *body.Box.Y + *body.Box.Height - tly
				}
				tly = tly - *body.Box.Y
			}
			// box height too small
			if height < 50 {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}

			// set the new relative reference to the newly cropped image
			tag.Origin.ImageSizeID = imageSizeID
			tag.Origin.Box.X = &tlx
			tag.Origin.Box.Y = &tly
			tag.Origin.Box.Width = &width
			tag.Origin.Box.Height = &height
		}
		i++
	}
	return imageFound, nil
}

func replaceImage(collection *mongo.Collection, imageReplace *types.Image, imageFile []byte) (*int64, error) {
	// replace or create the file
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := filepath.Join(folderDir, imageReplace.Origin, imageReplace.Name)
	err := os.WriteFile(path, imageFile, 0644)
	if err != nil {
		return nil, fmt.Errorf("os.WriteFile has failed: %v", err)
	}

	// get the new dimensions
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Open has failed: %v", err)
	}
	imageDecoded, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("image.DecodeConfig has failed. Only jpeg/jpg and png supported: %v", err)
	}

	// update in db the new dimensions, tag boxes and new size
	query := bson.M{"_id": imageReplace.ID}
	update := bson.M{
		"$set": bson.M{
			"width":  imageDecoded.Width,
			"height": imageDecoded.Height,
			"tags":   imageReplace.Tags,
			"size":   imageReplace.Size,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return nil, fmt.Errorf("UpdateOne has failed: %v", err)
	}
	return &res.ModifiedCount, nil
}

// InsertImageUnwanted insert an unwanted image
func InsertImageUnwanted(mongoClient *mongo.Client, body types.Image) (interface{}, error) {
	now := time.Now()
	body.CreationDate = &now
	body.Origin = strings.ToLower(body.Origin)

	// insert the unwanted image
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
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
