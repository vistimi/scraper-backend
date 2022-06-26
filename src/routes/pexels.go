package routes

import (
	"fmt"
	"scrapper/src/mongodb"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"os"
	"sort"

	"encoding/json"
	"path/filepath"
)

func SearchPhotosPexels(mongoClient *mongo.Client) (interface{}, error) {
	var insertedIds []primitive.ObjectID

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	origin := "pexels"
	err := os.MkdirAll(filepath.Join(folderDir, origin), os.ModePerm)
	if err != nil {
		return nil, err
	}

	// collectionPexels := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("PEXELS_COLLECTION"))

	unwantedTags, err := mongodb.TagsUnwantedNames(mongoClient)
	if err != nil {
		return nil, err
	}
	sort.Strings(unwantedTags)

	wantedTags, err := mongodb.TagsWantedNames(mongoClient)
	if err != nil {
		return nil, err
	}
	sort.Strings(wantedTags)

	for _, wantedTag := range wantedTags {
		page := 1
		searchPerPage, err := searchPhotosPerPagePexels(wantedTag, page)
		if err != nil {
			return nil, err
		}

		for page := page; page <= searchPerPage.TotalResults / searchPerPage.PerPage; page++ {
			searchPerPage, err = searchPhotosPerPagePexels(wantedTag, page)
			if err != nil {
				return nil, err
			}

			for _, photo := range searchPerPage.Photos {
			}
		}
		return searchPerPage, nil
	}
	return insertedIds, nil
}

type PhotoPexels struct {
	ID              int          `json:"id"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	URL             string       `json:"url"`
	Photographer    string       `json:"photographer"`
	PhotographerURL string       `json:"photographer_url"`
	PhotographerID  int          `json:"photographer_id"`
	AvgColor        string       `json:"avg_color"`
	Liked           bool         `json:"liked"`
	Src             SourcePexels `json:"src"`
}

type SourcePexels struct {
	Original  string `json:"original"`
	Large2X   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type SearchPhotoResponsePexels struct {
	TotalResults int            `json:"total_results"`
	Page         int            `json:"page"`
	PerPage      int            `json:"per_page"`
	Photos       []*PhotoPexels `json:"photos"`
	NextPage     string         `json:"next_page"`
	PrevPage     string         `json:"prev_page"`
}

func searchPhotosPerPagePexels(tag string, page int) (*SearchPhotoResponsePexels, error) {
	r := &Request{
		Host: "https://api.pexels.com/v1/search?",
		Args: map[string]string{
			"query": tag,
			"page": fmt.Sprint(page),
		},
		Header: map[string][]string{
			"Authorization": {utils.DotEnvVariable("PEXELS_PUBLIC_KEY")},
		},
	}
	fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var searchPerPage SearchPhotoResponsePexels
	err = json.Unmarshal(body, &searchPerPage)
	if err != nil {
		return nil, err
	}
	return &searchPerPage, nil
}
