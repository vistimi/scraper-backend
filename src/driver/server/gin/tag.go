package gin

import "context"

func (d DriverServerGin) CreateTagWanted(ctx context.Context, body ParamsReadPictureFile) (*DataSchema, error) {
	buffer, err := d.ControllerPicture.ReadPictureFile(ctx, params.Origin, params.Name, params.Extension)
	if err != nil {
		return nil, err
	}
	data := DataSchema{DataType: params.Extension, DataFile: buffer}
	return &data, nil
}

type ParamsRemoveTag struct {
	ID string `uri:"id" binding:"required"`
}

// func RemoveTagWanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
// 	collectionTagsWanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("TAGS_DESIRED_COLLECTION"))
// 	tagID, err := primitive.ObjectIDFromHex(params.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return mongodb.RemoveTag(collectionTagsWanted, tagID)
// }

// func RemoveTagUnwanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
// 	collectionTagsUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("TAGS_UNDESIRED_COLLECTION"))
// 	tagID, err := primitive.ObjectIDFromHex(params.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return mongodb.RemoveTag(collectionTagsUnwanted, tagID)
// }
