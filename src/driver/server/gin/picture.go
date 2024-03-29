package gin

import (
	"context"
	"fmt"

	"scraper-backend/src/driver/model"
	serverModel "scraper-backend/src/driver/server/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	Collection string `uri:"collection" binding:"required"`
	Origin     string `uri:"origin" binding:"required"`
}

// TODO: use Query and not Scan
func (d DriverServerGin) ReadPicturesID(ctx context.Context, params ParamsReadPicturesID) ([]serverModel.Picture, error) {
	// projEx := expression.NamesList(expression.Name("ID"))
	filtEx := expression.Name("Origin").Contains(params.Origin)
	controllerPictures, err := d.ControllerPicture.ReadPictures(ctx, params.Collection, nil, &filtEx)
	if err != nil {
		return nil, err
	}
	driverServerPictures := make([]serverModel.Picture, 0, len(controllerPictures))
	for _, controllerPicture := range controllerPictures {
		var serverPicture serverModel.Picture
		serverPicture.DriverMarshal(controllerPicture)
		driverServerPictures = append(driverServerPictures, serverPicture)
	}
	return driverServerPictures, nil
}

type ParamsReadPicture struct {
	Origin     string `uri:"origin" binding:"required"`
	ID         string `uri:"id" binding:"required"`
	Collection string `uri:"collection" binding:"required"`
}

func (d DriverServerGin) ReadPicture(ctx context.Context, params ParamsReadPicture) (*serverModel.Picture, error) {
	id, err := model.ParseUUID(params.ID)
	if err != nil {
		return nil, err
	}
	controllerPicture, err := d.ControllerPicture.ReadPicture(ctx, params.Collection, params.Origin, id)
	if err != nil || controllerPicture == nil {
		return nil, err
	}
	var driverServerPicture serverModel.Picture
	driverServerPicture.DriverMarshal(*controllerPicture)
	return &driverServerPicture, nil
}

func (d DriverServerGin) ReadPicturesBlocked(ctx context.Context) ([]serverModel.Picture, error) {
	controllerPictures, err := d.ControllerPicture.ReadPictures(ctx, "blocked", nil, nil)
	if err != nil {
		return nil, err
	}
	driverServerPictures := make([]serverModel.Picture, 0, len(controllerPictures))
	for _, controllerPicture := range controllerPictures {
		var serverPicture serverModel.Picture
		serverPicture.DriverMarshal(controllerPicture)
		driverServerPictures = append(driverServerPictures, serverPicture)
	}
	return driverServerPictures, nil
}

type ParamsDeletePictureAndFile struct {
	Origin string `uri:"origin" binding:"required"`
	ID     string `uri:"id" binding:"required"`
	Name   string `uri:"name" binding:"required"`
}

func (d DriverServerGin) DeletePictureAndFile(ctx context.Context, params ParamsDeletePictureAndFile) (string, error) {
	id, err := model.ParseUUID(params.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.DeletePictureAndFile(ctx, params.Origin, id, params.Name); err != nil {
		return "error", err
	}
	return "ok", nil
}

type ParamsDeletePicture struct {
	Origin string `uri:"origin" binding:"required"`
	ID     string `uri:"id" binding:"required"`
}

func (d DriverServerGin) DeletePicture(ctx context.Context, params ParamsDeletePicture) (string, error) {
	id, err := model.ParseUUID(params.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.DeletePicture(ctx, params.Origin, id); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyUpdatePictureTag struct {
	Origin *string                 `json:"origin"`
	ID     *string                 `json:"id"`
	Tag    *serverModel.PictureTag `json:"tag"`
}

func (d DriverServerGin) UpdatePictureTag(ctx context.Context, body BodyUpdatePictureTag) (string, error) {
	if body.Origin == nil || body.ID == nil || body.Tag == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	if !body.Tag.BoxInformation.IsValid() {
		return "error", fmt.Errorf("body not valid, tag.boxInformation missing")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.UpdatePictureTag(ctx, *body.Origin, id, model.NewUUID(), body.Tag.DriverUnmarshal()); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyDeletePictureTag struct {
	Origin *string `json:"origin"`
	ID     *string `json:"id"`
	TagID  *string `json:"tagID"`
}

func (d DriverServerGin) DeletePictureTag(ctx context.Context, body BodyDeletePictureTag) (string, error) {
	if body.Origin == nil || body.ID == nil || body.TagID == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	tagID, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.DeletePictureTag(ctx, *body.Origin, id, tagID); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyUpdatePictureCrop struct {
	Origin *string          `json:"origin"`
	ID     *string          `json:"id"`
	Name   *string          `json:"name"`
	Box    *serverModel.Box `json:"box"`
}

func (d DriverServerGin) UpdatePictureCrop(ctx context.Context, body BodyUpdatePictureCrop) (string, error) {
	if body.Origin == nil || body.ID == nil || body.Name == nil || body.Box == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.UpdatePictureCrop(ctx, *body.Origin, id, *body.Name, model.NewUUID(), body.Box.DriverUnmarshal()); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyCreatePictureCrop struct {
	Origin        *string          `json:"origin"`
	ID            *string          `json:"id"`
	Name          *string          `json:"name"`
	PictureSizeID *string          `json:"pictureSizeID"`
	Box           *serverModel.Box `json:"box"`
}

func (d DriverServerGin) CreatePictureCrop(ctx context.Context, body BodyCreatePictureCrop) (string, error) {
	if body.Origin == nil || body.ID == nil || body.Name == nil || body.PictureSizeID == nil || body.Box == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	pictureSizeID, err := model.ParseUUID(*body.PictureSizeID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.CreatePictureCrop(ctx, *body.Origin, id, model.NewUUID(), pictureSizeID, body.Box.DriverUnmarshal()); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyCreatePictureCopy struct {
	Origin *string `json:"origin"`
	ID     *string `json:"id"`
}

func (d DriverServerGin) CreatePictureCopy(ctx context.Context, body BodyCreatePictureCopy) (string, error) {
	if body.Origin == nil || body.ID == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.CreatePictureCopy(ctx, *body.Origin, id, model.NewUUID()); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyUpdatePictureTransfer struct {
	Origin *string `json:"origin"`
	ID     *string `json:"id"`
	From   *string `json:"from"`
	To     *string `json:"to"`
}

func (d DriverServerGin) UpdatePictureTransfer(ctx context.Context, body BodyUpdatePictureTransfer) (string, error) {
	if body.Origin == nil || body.ID == nil || body.From == nil || body.To == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.UpdatePictureTransfer(ctx, *body.Origin, id, *body.From, *body.To); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyCreatePictureBlocked struct {
	Origin *string `json:"origin"`
	ID     *string `json:"id"`
}

func (d DriverServerGin) CreatePictureBlocked(ctx context.Context, body BodyCreatePictureBlocked) (string, error) {
	if body.Origin == nil || body.ID == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.UpdatePictureTransfer(ctx, *body.Origin, id, "process", "blocked"); err != nil {
		return "error", err
	}
	return "ok", nil
}

type BodyDeletePictureBlocked struct {
	Origin *string `json:"origin"`
	ID     *string `json:"id"`
}

func (d DriverServerGin) DeletePictureBlocked(ctx context.Context, body BodyCreatePictureBlocked) (string, error) {
	if body.Origin == nil || body.ID == nil {
		return "error", fmt.Errorf("body fields must not be empty")
	}
	id, err := model.ParseUUID(*body.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerPicture.UpdatePictureTransfer(ctx, *body.Origin, id, "blocked", "process"); err != nil {
		return "error", err
	}
	return "ok", nil
}
