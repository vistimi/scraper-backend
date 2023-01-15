package gin

import (
	"context"

	driverServerModel "scraper-backend/src/driver/server/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
)

type ParamsReadPictureFile struct {
	Origin    string `uri:"origin" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Extension string `uri:"extension" binding:"required"`
}

func (d DriverServerGin) ReadPictureFile(ctx context.Context, params ParamsReadPictureFile) (*DataSchema, error) {
	buffer, err := d.ControllerPicture.ReadPictureFile(ctx, params.Origin, params.Name, params.Extension)
	if err != nil {
		return nil, err
	}
	data := DataSchema{DataType: params.Extension, DataFile: buffer}
	return &data, nil
}

type ParamsReadPicturesID struct {
	Origin     string `uri:"origin" binding:"required"`
	Collection string `uri:"collection" binding:"required"`
}

// FindImagesIDs get all the IDs of a picture
func (d DriverServerGin) ReadPicturesID(ctx context.Context, params ParamsReadPicturesID) ([]driverServerModel.Picture, error) {
	projEx := expression.NamesList(expression.Name("ID"))
	return d.ControllerPicture.ReadPictures(ctx, params.Collection, &projEx, nil)
}

type ParamsReadPicture struct {
	Origin     string    `uri:"origin" binding:"required"`
	ID         uuid.UUID `uri:"id" binding:"required"`
	Collection string    `uri:"collection" binding:"required"`
}

// FindImage get a specific image
func (d DriverServerGin) ReadPicture(ctx context.Context, params ParamsReadPicture) (*driverServerModel.Picture, error) {
	return d.ControllerPicture.ReadPicture(ctx, params.Collection, params.Origin, params.ID)
}

// FindImagesUnwanted get all the unwanted images
func (d DriverServerGin) FindPicturesBlocked(ctx context.Context) ([]driverServerModel.Picture, error) {
	return d.ControllerPicture.ReadPictures(ctx, "blocked", nil, nil)
}

type ParamsDeletePictureAndFile struct {
	Origin string    `uri:"origin" binding:"required"`
	ID     uuid.UUID `uri:"id" binding:"required"`
	Name   string    `uri:"name" binding:"required"`
}

// RemoveImageAndFile removes in db and file of a pending image
func (d DriverServerGin) DeletePictureAndFile(ctx context.Context, params ParamsDeletePictureAndFile) (string, error) {
	if err := d.ControllerPicture.DeletePictureAndFile(ctx, params.Origin, params.ID, params.Name); err != nil {
		return "error", err
	}
	return "ok", nil
}

type ParamsDeletePicture struct {
	Origin string    `uri:"origin" binding:"required"`
	ID     uuid.UUID `uri:"id" binding:"required"`
}

// RemoveImage removes in db an unwanted image
func (d DriverServerGin) DeletePicture(ctx context.Context, params ParamsDeletePicture) (string, error) {
	if err := d.ControllerPicture.DeletePicture(ctx, params.Origin, params.ID); err != nil {
		return "error", err
	}
	return "ok", nil
}

// type BodyUpdatePictureTag struct {
// 	Origin string                     `json:"origin,omitempty"`
// 	ID     uuid.UUID                  `json:"id,omitempty"`
// 	Tag    controllerModel.PictureTag `json:"tag,omitempty"`
// }

// // UpdateImageTagsPush add tags to a pending image
// func (d DriverServerGin) UpdatePictureTag(ctx context.Context, body BodyUpdatePictureTag) (*int64, error) {
// 	if body.Tag.Origin.Box.Tlx == nil || tag.Origin.Box.Tly == nil || tag.Origin.Box.Width == nil || tag.Origin.Box.Height == nil {
// 		return nil, fmt.Errorf("body not valid, box fields missing: %v", tag.Origin.Box)
// 	}
// 	d.ControllerPicture.UpdatePictureTag(ctx, body.Origin, body.ID, uuid.New(), body.Tag)
// 	if body.ID == primitive.NilObjectID {
// 		return nil, errors.New("body not valid, ID empty")
// 	}
// 	for _, tag := range body.Tags {

// 	}
// 	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))
// 	return mongodb.UpdateImageTagsPush(collectionImagesPending, body)
// }

// // UpdateImageTagsPush remove tags to a pending image
// func (u *usecasePictures) UpdateImageTagsPull(mongoClient *mongo.Client, body types.BodyUpdateImageTagsPull) (*int64, error) {
// 	if body.ID == primitive.NilObjectID {
// 		return nil, errors.New("body not valid, ID empty")
// 	}
// 	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))
// 	return mongodb.UpdateImageTagsPull(collectionImagesPending, body)
// }
