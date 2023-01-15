package controller

import (
	"context"
	"fmt"
	"regexp"
	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceDatabase "scraper-backend/src/driver/interface/database"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type ControllerTag struct {
	Dynamodb          interfaceDatabase.DriverDynamodbTag
	ControllerPicture ControllerPicture
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

// func (c ControllerTag) ReadTag(ctx context.Context, tag controllerModel.Tag) (*controllerModel.Tag, error) {
// 	return c.Dynamodb.ReadTag(ctx, tag.Type, tag.Name)
// }

func (c ControllerTag) ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error) {
	return c.Dynamodb.ReadTags(ctx, primaryKey)
}

// // TagsWanted find all the names of wanted tags
// func TagsWantedNames(mongoClient *mongo.Client) ([]string, error) {
// 	collectionTagsWanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("TAGS_DESIRED_COLLECTION"))
// 	res, err := FindMany[types.Tag](collectionTagsWanted, bson.M{})
// 	if err != nil {
// 		return nil, fmt.Errorf("FindTags Wanted has failed: \n%v", err)
// 	}
// 	var wantedTags []string
// 	for _, tag := range res {
// 		wantedTags = append(wantedTags, strings.ToLower(tag.Name))
// 	}
// 	return wantedTags, nil
// }

// // TagsUnwantednames find all the names of wanted tags
// func TagsUnwantedNames(mongoClient *mongo.Client) ([]string, error) {
// 	collectionTagsUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("TAGS_UNDESIRED_COLLECTION"))
// 	res, err := FindMany[types.Tag](collectionTagsUnwanted, bson.M{})
// 	if err != nil {
// 		return nil, fmt.Errorf("FindTags Unwated has failed: \n%v", err)
// 	}
// 	var unwantedTags []string
// 	for _, tag := range res {
// 		unwantedTags = append(unwantedTags, strings.ToLower(tag.Name))
// 	}
// 	return unwantedTags, nil
// }

// func TagsNames(mongoClient *mongo.Client) ([]string, []string, error) {
// 	unwantedTags, err := TagsUnwantedNames(mongoClient)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	if (unwantedTags == nil) || (len(unwantedTags) == 0) {
// 		return nil, nil, errors.New("unwantedTags are empty")
// 	}
// 	sort.Strings(unwantedTags)

// 	wantedTags, err := TagsWantedNames(mongoClient)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	if (wantedTags == nil) || (len(wantedTags) == 0) {
// 		return nil, nil, errors.New("wantedTags are empty")
// 	}
// 	sort.Strings(wantedTags)
// 	return unwantedTags, wantedTags, nil
// }
