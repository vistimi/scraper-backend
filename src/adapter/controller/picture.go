package controller

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceDatabase "scraper-backend/src/driver/interface/database"
	interfaceStorage "scraper-backend/src/driver/interface/storage"
	model "scraper-backend/src/driver/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type ControllerPicture struct {
	S3                 interfaceStorage.DriverS3
	BucketName         string
	DynamodbProcess    interfaceDatabase.DriverDynamodbPicture
	DynamodbValidation interfaceDatabase.DriverDynamodbPicture
	DynamodbProduction interfaceDatabase.DriverDynamodbPicture
	DynamodbBlocked    interfaceDatabase.DriverDynamodbPicture
}

func (c ControllerPicture) driverDynamodbMap(state string) (interfaceDatabase.DriverDynamodbPicture, error) {
	switch state {
	case "production":
		return c.DynamodbProduction, nil
	case "validation":
		return c.DynamodbValidation, nil
	case "process":
		return c.DynamodbProcess, nil
	case "blocked":
		return c.DynamodbProcess, nil
	default:
		return nil, fmt.Errorf("table name %s not available", state)
	}
}

func (c ControllerPicture) ReadPictures(ctx context.Context, state string, projection *expression.ProjectionBuilder, filter *expression.ConditionBuilder) ([]controllerModel.Picture, error) {
	dynamodb, err := c.driverDynamodbMap(state)
	if err != nil {
		return nil, err
	}
	return dynamodb.ReadPictures(ctx, projection, filter)
}

func (c ControllerPicture) ReadPicture(ctx context.Context, state string, primaryKey string, sortKey model.UUID) (*controllerModel.Picture, error) {
	dynamodb, err := c.driverDynamodbMap(state)
	if err != nil {
		return nil, err
	}
	return dynamodb.ReadPicture(ctx, primaryKey, sortKey)
}

func (c ControllerPicture) ReadPictureFile(ctx context.Context, origin, name, extension string) ([]byte, error) {
	fileName := fmt.Sprintf("%s.%s", name, extension)
	path := filepath.Join(origin, fileName)

	buffer, err := c.S3.ItemRead(ctx, c.BucketName, path)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (c ControllerPicture) CreatePicture(ctx context.Context, id model.UUID, picture controllerModel.Picture, buffer []byte) error {
	path := fmt.Sprintf("%s/%s.%s", picture.Origin, picture.Name, picture.Extension)
	if err := c.S3.ItemCreate(ctx, bytes.NewReader(buffer), c.BucketName, path); err != nil {
		return err
	}
	return c.DynamodbProcess.CreatePicture(ctx, id, picture)
}

func (c ControllerPicture) DeletePicture(ctx context.Context, primaryKey string, sortKey model.UUID) error {
	return c.DynamodbProcess.DeletePicture(ctx, primaryKey, sortKey)
}

func (c ControllerPicture) DeletePictureAndFile(ctx context.Context, primaryKey string, sortKey model.UUID, name string) error {
	if err := c.DynamodbProcess.DeletePicture(ctx, primaryKey, sortKey); err != nil {
		return err
	}
	path := filepath.Join(primaryKey, name)
	return c.S3.ItemDelete(ctx, c.BucketName, path)
}

func (c ControllerPicture) DeletePicturesAndFiles(ctx context.Context, pictures []controllerModel.Picture) error {
	for _, picture := range pictures {
		if err := c.DynamodbProcess.DeletePicture(ctx, picture.Origin, picture.ID); err != nil {
			return err
		}
		path := filepath.Join(picture.Origin, picture.Name)
		if err := c.S3.ItemDelete(ctx, c.BucketName, path); err != nil {
			return err
		}
	}
	return nil
}

func (c ControllerPicture) CreatePictureTag(ctx context.Context, primaryKey string, sortKey model.UUID, tagID model.UUID, tag controllerModel.PictureTag) error {
	if err := c.DynamodbProcess.CreatePictureTag(ctx, primaryKey, sortKey, tagID, tag); err != nil {
		return err
	}
	return nil
}

func (c ControllerPicture) UpdatePictureTag(ctx context.Context, primaryKey string, sortKey model.UUID, tagID model.UUID, tag controllerModel.PictureTag) error {
	if err := c.DynamodbProcess.UpdatePictureTag(ctx, primaryKey, sortKey, tagID, tag); err != nil {
		return err
	}
	return nil
}

func (c ControllerPicture) DeletePictureTag(ctx context.Context, primaryKey string, sortKey model.UUID, tagID model.UUID) error {
	if err := c.DynamodbProcess.DeletePictureTag(ctx, primaryKey, sortKey, tagID); err != nil {
		return err
	}
	return nil
}

func (c ControllerPicture) UpdatePictureCrop(ctx context.Context, primaryKey string, sortKey model.UUID, name string, pictureSizeID model.UUID, box controllerModel.Box) error {
	newPicture, err := c.cropPicture(ctx, box, primaryKey, sortKey, pictureSizeID)
	if err != nil {
		return err
	}

	newFile, err := c.cropFile(ctx, box, primaryKey, name)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s.%s", newPicture.Origin, newPicture.Name, newPicture.Extension)
	buffer, err := fileToBuffer(*newPicture, newFile)
	if err != nil {
		return err
	}
	if err := c.S3.ItemCreate(ctx, buffer, c.BucketName, path); err != nil {
		return err
	}

	// update current image
	if err := c.DynamodbProcess.CreatePicture(ctx, newPicture.ID, *newPicture); err != nil {
		return err
	}
	return nil
}

func (c ControllerPicture) CreatePictureCrop(ctx context.Context, primaryKey string, sortKey model.UUID, id model.UUID, pictureSizeID model.UUID, box controllerModel.Box) error {
	newPicture, err := c.DynamodbProcess.ReadPicture(ctx, primaryKey, sortKey)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s_%s", newPicture.OriginID, time.Now().Format(time.RFC3339))
	newPicture.Name = name
	newPicture.CreationDate = time.Now()

	newPicture, err = c.cropPicture(ctx, box, primaryKey, sortKey, pictureSizeID)
	if err != nil {
		return err
	}

	newFile, err := c.cropFile(ctx, box, primaryKey, name)
	if err != nil {
		return err
	}

	destinationPath := fmt.Sprintf("%s/%s.%s", newPicture.Origin, name, newPicture.Extension)
	buffer, err := fileToBuffer(*newPicture, newFile)
	if err != nil {
		return err
	}
	if err := c.S3.ItemCreate(ctx, buffer, c.BucketName, destinationPath); err != nil {
		return err
	}

	if err := c.DynamodbProcess.CreatePicture(ctx, id, *newPicture); err != nil {
		return err
	}
	return nil
}

func (c ControllerPicture) CreatePictureCopy(ctx context.Context, primaryKey string, sortKey model.UUID, id model.UUID) error {
	newPicture, err := c.DynamodbProcess.ReadPicture(ctx, primaryKey, sortKey)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s_%s", newPicture.OriginID, time.Now().Format(time.RFC3339))
	newPicture.Name = name
	newPicture.CreationDate = time.Now()

	sourcePath := fmt.Sprintf("%s/%s.%s", newPicture.Origin, newPicture.Name, newPicture.Extension)
	destinationPath := fmt.Sprintf("%s/%s.%s", newPicture.Origin, name, newPicture.Extension)
	if err := c.S3.ItemCopy(ctx, c.BucketName, sourcePath, destinationPath); err != nil {
		return err
	}

	if err := c.DynamodbProcess.CreatePicture(ctx, id, *newPicture); err != nil {
		return err
	}
	return nil
}

func (c ControllerPicture) UpdatePictureTransfer(ctx context.Context, primaryKey string, sortKey model.UUID, from, to string) error {
	fromDynamodb, err := c.driverDynamodbMap(from)
	if err != nil {
		return err
	}

	oldPicture, err := fromDynamodb.ReadPicture(ctx, primaryKey, sortKey)
	if err != nil {
		return err
	}

	toDynamodb, err := c.driverDynamodbMap(to)
	if err != nil {
		return err
	}

	if err := toDynamodb.CreatePicture(ctx, oldPicture.ID, *oldPicture); err != nil {
		return err
	}

	if err := fromDynamodb.DeletePicture(ctx, primaryKey, sortKey); err != nil {
		return err
	}

	return nil
}

func (c ControllerPicture) CreatePictureBlocked(ctx context.Context, primaryKey string, sortKey model.UUID) error {
	picture, err := c.DynamodbProcess.ReadPicture(ctx, primaryKey, sortKey)
	if err != nil {
		return err
	}

	sourcePath := fmt.Sprintf("%s/%s.%s", picture.Origin, picture.Name, picture.Extension)
	if err := c.S3.ItemDelete(ctx, c.BucketName, sourcePath); err != nil {
		return err
	}

	picture.CreationDate = time.Now()
	if err := c.DynamodbBlocked.CreatePicture(ctx, picture.ID, *picture); err != nil {
		return err
	}

	if err := c.DynamodbProcess.DeletePicture(ctx, picture.Origin, picture.ID); err != nil {
		return err
	}

	return nil
}

func (c ControllerPicture) DeletePictureBlocked(ctx context.Context, primaryKey string, sortKey model.UUID) error {
	return c.DynamodbBlocked.DeletePicture(ctx, primaryKey, sortKey)
}

func fileToBuffer(picture controllerModel.Picture, file image.Image) (*bytes.Buffer, error) {
	// create buffer
	buffer := new(bytes.Buffer)
	// encode image to buffer

	if picture.Extension == "jpeg" || picture.Extension == "jpg" {
		err := jpeg.Encode(buffer, file, nil)
		if err != nil {
			return nil, fmt.Errorf("jpeg.Encode has failed: %v", err)
		}
	} else if picture.Extension == "png" {
		err := png.Encode(buffer, file)
		if err != nil {
			return nil, fmt.Errorf("png.Encode has failed: %v", err)
		}
	} else {
		return nil, fmt.Errorf("no image extension matching the buffer conversion")
	}
	return buffer, nil
}

func (c ControllerPicture) cropPicture(ctx context.Context, box controllerModel.Box, primaryKey string, sortKey model.UUID, pictureSizeID model.UUID) (*controllerModel.Picture, error) {
	oldPicture, err := c.DynamodbProcess.ReadPicture(ctx, primaryKey, sortKey)
	if err != nil {
		return nil, err
	}
	return updatePictureTagBoxes(box, *oldPicture, pictureSizeID)
}

func updatePictureTagBoxes(box controllerModel.Box, picture controllerModel.Picture, pictureSizeID model.UUID) (*controllerModel.Picture, error) {
	// new size creation
	size := controllerModel.PictureSize{
		ID:           model.NewUUID(),
		CreationDate: time.Now(),
		Box:          box, // absolute position
	}
	picture.Sizes = append(picture.Sizes, size)

	i := 0
	for {
		if i >= len(picture.Tags) {
			break
		}
		tag := &picture.Tags[i]

		if tag.BoxInformation.Valid {
			boxInformation := tag.BoxInformation.Body
			// relative position of tags
			tlx := boxInformation.Box.Tlx
			tly := boxInformation.Box.Tly
			width := boxInformation.Box.Width
			height := boxInformation.Box.Height

			// box outside on the image right
			if tlx > box.Tlx+box.Width {
				removeSliceElement(picture.Tags, i)
				continue
			}
			// box left outside on the image left
			if tlx < box.Tlx {
				if tlx+width < box.Tlx {
					// box outside on the image left
					width = 0
				} else {
					// box right inside the image
					width = width - box.Tlx + tlx
				}
				tlx = box.Tlx
			} else { // box left inside image
				if tlx+width > box.Tlx+box.Width {
					// box right outside on the image right
					width = box.Tlx + box.Width - tlx
				}
				tlx = tlx - box.Tlx
			}
			// box width too small
			if width < 50 {
				removeSliceElement(picture.Tags, i)
				continue
			}

			// box outside at the image bottom
			if tly > box.Tly+box.Height {
				removeSliceElement(picture.Tags, i)
				continue
			}
			// box top outside on the image top
			if tly < box.Tly {
				if tly+height < box.Tly {
					// box outside on the image top
					height = 0
				} else {
					// box bottom inside the image
					height = height - box.Tly + tly
				}
				tly = box.Tly
			} else { // box top inside image
				// box bottom outside on the image bottom
				if tly+height > box.Tly+box.Height {
					height = box.Tly + box.Height - tly
				}
				tly = tly - box.Tly
			}
			// box height too small
			if height < 50 {
				removeSliceElement(picture.Tags, i)
				continue
			}

			// set the new relative reference to the newly cropped image
			tag.BoxInformation.Body.PictureSizeID = pictureSizeID
			tag.BoxInformation.Body.Box.Tlx = tlx
			tag.BoxInformation.Body.Box.Tly = tly
			tag.BoxInformation.Body.Box.Width = width
			tag.BoxInformation.Body.Box.Height = height
		}
		i++
	}
	return &picture, nil
}

func (c ControllerPicture) cropFile(ctx context.Context, box controllerModel.Box, primaryKey string, name string) (image.Image, error) {
	path := filepath.Join(primaryKey, name)
	buffer, err := c.S3.ItemRead(ctx, c.BucketName, path)
	if err != nil {
		return nil, err
	}

	// convert []byte to image
	img, _, _ := image.Decode(bytes.NewReader(buffer))

	// crop the image with the bounding box rectangle
	cropRect := image.Rect(box.Tlx, box.Tly, box.Tlx+box.Width, box.Tly+box.Height)
	img, err = updateFileDimension(img, cropRect)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func updateFileDimension(img image.Image, cropRect image.Rectangle) (image.Image, error) {
	//Interface for asserting whether `img`
	//implements SubImage or not.
	//This can be defined globally.
	type CropableImage interface {
		image.Image
		SubImage(r image.Rectangle) image.Image
	}

	if p, ok := img.(CropableImage); ok {
		// Call SubImage. This should be fast,
		// since SubImage (usually) shares underlying pixel.
		return p.SubImage(cropRect), nil
	} else if cropRect = cropRect.Intersect(img.Bounds()); !cropRect.Empty() {
		// If `img` does not implement `SubImage`,
		// copy (and silently convert) the image portion to RGBA image.
		rgbaImg := image.NewRGBA(cropRect)
		for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
			for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
		return rgbaImg, nil
	} else {
		return nil, fmt.Errorf("cannot crop the image")
	}
}

func removeSliceElement[T any](slice []T, i int) {
	if i == len(slice)-1 {
		// last element removed
		slice = slice[:i]
	} else {
		// not last element removed
		slice = append(slice[:i], slice[i+1:]...)
	}
}
