package routes

import (
	"fmt"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hbagdi/go-unsplash/unsplash"

	"encoding/json"
)

type ParamsSearchPhotoUnsplash struct {
	Quality string `uri:"quality" binding:"required"`
}

func SearchPhotosUnsplash(mongoClient *mongo.Client, params ParamsSearchPhotoUnsplash) (interface{}, error) {
	r := &Request{
		Host: "https://api.unsplash.com/photos/?",
		Args: map[string]string{
			"client_id": utils.DotEnvVariable("UNSPLASH_PUBLIC_KEY"),
			"page": "1",
		},
	}

	fmt.Println(r.URL())

	body, err := r.Execute()
	if err != nil {
		return nil, err
	}

	photos := make([]unsplash.Photo, 0)
	err = json.Unmarshal(body, &photos)
	if err != nil {
		return nil, err
	}
	return &photos, nil
}
