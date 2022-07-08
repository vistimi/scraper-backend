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
func RemoveImage(collection *mongo.Collection, id primitive.ObjectID, origin string) (*int64, error) {
	query := bson.M{"_id": id, "origin": origin}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// RemoveImageAndFile remove an image based on its mongodb id and remove its file
func RemoveImageAndFile(collection *mongo.Collection, id primitive.ObjectID, origin string) (*int64, error) {
	image, err := FindOne[types.Image](collection, bson.M{"_id": id, "origin": origin})
	if err != nil {
		return nil, fmt.Errorf("FindImageByID has failed: %v", err)
	}
	deletedCount, err := RemoveImage(collection, id, origin)
	if err != nil {
		return nil, fmt.Errorf("RemoveImage has failed: %v", err)
	}
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := fmt.Sprintf(filepath.Join(folderDir, origin, image.Name))
	err = os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("os.Remove has failed: %v", err)
	}
	return deletedCount, nil
}

func RemoveImagesAndFilesOneOrigin(mongoClient *mongo.Client, origin string, query bson.M, options *options.FindOptions) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	var deletedCount int64
	images, err := FindMany[types.Image](collectionImages, query, options)
	if err != nil {
		return nil, fmt.Errorf("FindImagesIDs has failed: %v", err)
	}
	for _, image := range images {
		deletedOne, err := RemoveImageAndFile(collectionImages, image.ID, origin)
		if err != nil {
			return nil, fmt.Errorf("RemoveImageAndFile has failed: %v", err)
		}
		deletedCount += *deletedOne
	}
	return &deletedCount, nil
}

// Remove all the images in DB and their related file matching the query and options given, for all origins
func RemoveImagesAndFilesAllOrigins(mongoClient *mongo.Client, query bson.M, options *options.FindOptions) (*int64, error) {

	imageOrigins := utils.ImageOrigins()
	var deletedCount int64
	for _, origin := range imageOrigins {
		count, err := RemoveImagesAndFilesOneOrigin(mongoClient, origin, query, options)
		if err != nil {
			return nil, fmt.Errorf("RemoveImageAndFile has failed: %v", err)
		}
		deletedCount += *count
	}
	return &deletedCount, nil
}

// UpdateImageTags add tags to an image based on its mongodb id
func UpdateImageTagsPush(collection *mongo.Collection, body types.BodyUpdateImageTagsPush) (*types.Image, error) {
	query := bson.M{"_id": body.ID}
	if body.Tags != nil {
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
		_, err := collection.UpdateOne(context.TODO(), query, update)
		if err != nil {
			return nil, fmt.Errorf("UpdateOne has failed: %v", err)
		}
	}
	return FindOne[types.Image](collection, query)
}

func UpdateImageTagsPull(collection *mongo.Collection, body types.BodyUpdateImageTagsPull) (interface{}, error) {
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

func UpdateImageFile(collection *mongo.Collection, body types.BodyUpdateImageFile) (*types.Image, error) {
	imageSizeID := primitive.NewObjectID()
	now := time.Now()
	size := types.ImageSize{
		ID:           imageSizeID,
		CreationDate: &now,
		Box:          body.Box, // absolute position
	}

	query := bson.M{"origin": body.Origin, "name": body.Name}
	imageFound, err := FindOne[types.Image](collection, query)
	if err != nil {
		return nil, fmt.Errorf("FindOne has failed: %v", err)
	}
	i := 0
	fmt.Println(*body.Box.X, *body.Box.Y, *body.Box.Width, *body.Box.Height)
	for {
		if i >= len(imageFound.Tags) {
			break
		}

		tag := imageFound.Tags[i]
		fmt.Println(tag.Name)
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

	// replace the file
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := fmt.Sprintf(filepath.Join(folderDir, body.Origin, body.Name))
	err = os.WriteFile(path, body.File, 0644)
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
	update := bson.M{
		"$set": bson.M{
			"width":  imageDecoded.Width,
			"height": imageDecoded.Height,
			"tags":   imageFound.Tags,
		},
		"$push": bson.M{
			"size": size,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return nil, fmt.Errorf("UpdateOne has failed: %v", err)
	}
	return FindOne[types.Image](collection, query)
}

func InsertImageUnwanted(mongoClient *mongo.Client, body types.Image) (interface{}, error) {
	if body.Origin == "" || body.OriginID == "" {
		return nil, errors.New("Some fields are empty!")
	}
	now := time.Now()
	body.CreationDate = &now
	body.Origin = strings.ToLower(body.Origin)

	// insert the unwanted image
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	query := bson.M{"origin": body.Origin, "originID": body.OriginID}
	insertedID, err := InsertOne(collectionImagesUnwanted, body, query)
	if err != nil {
		return nil, fmt.Errorf("insertUser has failed: %v", err)
	}
	return insertedID, nil
}
