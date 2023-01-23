package controller

import (
	"context"
	"fmt"
	"regexp"
	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceDatabase "scraper-backend/src/driver/interface/database"
	interfaceAdapter "scraper-backend/src/adapter/interface"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type ControllerTag struct {
	Dynamodb          interfaceDatabase.DriverDynamodbTag
	ControllerPicture interfaceAdapter.ControllerPicture
}

func (c ControllerTag) CreateTag(ctx context.Context, tag controllerModel.Tag) error {
	existingTags, err := c.Dynamodb.ScanTags(ctx)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(existingTags, func(existingTag controllerModel.Tag) bool {
		regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, strings.ToLower(existingTag.Name))
		matched, err := regexp.Match(regexpMatch, []byte(tag.Name)) // e.g. match if thisTag has `model` and bodyTag `models`
		if err != nil {
			return false
		}
		return matched
	})
	if idx != -1 {
		return fmt.Errorf("tag `%s` is too closely related to `%+#v`", tag.Name, existingTags[idx])
	}

	tag.ID = uuid.New()
	tag.CreationDate = time.Now()
	tag.Name = strings.ToLower(tag.Name)
	tag.OriginName = strings.ToLower(tag.OriginName)
	return c.Dynamodb.CreateTag(ctx, tag)
}

func (c ControllerTag) CreateTagBlocked(ctx context.Context, tag controllerModel.Tag) error {
	if err := c.CreateTag(ctx, tag); err != nil {
		return err
	}

	projEx := expression.NamesList(expression.Name("Origin"), expression.Name("Name"), expression.Name("Tags"))
	filtEx := expression.Name("Tags.#Name").Contains(tag.Name)
	pictures, err := c.ControllerPicture.ReadPictures(ctx, "process", &projEx, &filtEx)
	if err != nil {
		return err
	}

	return c.ControllerPicture.DeletePicturesAndFiles(ctx, pictures)
}

func (c ControllerTag) DeleteTag(ctx context.Context, primaryKey string, sortKey uuid.UUID) error {
	return c.Dynamodb.DeleteTag(ctx, primaryKey, sortKey)
}

func (c ControllerTag) ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error) {
	return c.Dynamodb.ReadTags(ctx, primaryKey)
}
