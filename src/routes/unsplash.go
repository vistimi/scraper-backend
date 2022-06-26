package routes

import (
	"fmt"
	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hbagdi/go-unsplash/unsplash"

	"encoding/json"

	"os"
	"sort"

	"path/filepath"
	"strings"

	"strconv"
)

func SearchPhotosUnsplash(mongoClient *mongo.Client) ([]primitive.ObjectID, error) {
	var insertedIds []primitive.ObjectID

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	origin := "unsplash"
	err := os.MkdirAll(filepath.Join(folderDir, origin), os.ModePerm)
	if err != nil {
		return nil, err
	}

	collectionUnsplash := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNSPLASH_COLLECTION"))

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

		searchPerPage, err := searchPhotosPerPageUnsplash(wantedTag, page)
		if err != nil {
			return nil, err
		}

		for page := page; page <= int(*searchPerPage.TotalPages); page++ {
			searchPerPage, err = searchPhotosPerPageUnsplash(wantedTag, page)
			if err != nil {
				return nil, err
			}

			for _, photo := range *searchPerPage.Results {
				// extract the photo informations
				var photoTags []string
				for _, tag := range *photo.Tags {
					photoTags = append(photoTags, strings.ToLower(*tag.Title))
				}

				// skip image if one of its tag is unwanted
				idx := utils.FindIndexRegExp(unwantedTags, photoTags)
				if idx != -1 {
					continue // skip image with unwated tag
				}

				//find download link and extension
				link := photo.Urls.Small
				extension := link.Query().Get("fm")

				// download photo into folder and rename it <id>.<format>
				fileName := fmt.Sprintf("%s.%s", *photo.ID, extension)
				path := fmt.Sprintf(filepath.Join(folderDir, origin, fileName))
				err = DownloadFile(link.String(), path)
				if err != nil {
					return nil, err
				}

				// tags creation
				var tags []types.Tag
				now := time.Now()
				for _, photoTag := range *photo.Tags {
					tag := types.Tag{
						Name:         strings.ToLower(*photoTag.Title),
						Origin:       "unsplash",
						CreationDate: &now,
					}
					tags = append(tags, tag)
				}

				// image creation
				width, err := strconv.Atoi(link.Query().Get("w"))
				if err != nil {
					return nil, err
				}
				var height int
				if photo.Height != nil && photo.Width != nil {
					height = *photo.Height * width / *photo.Width
				}
				var originId string
				if photo.ID != nil {
					originId = *photo.ID
				}
				var title string
				if photo.Description != nil {
					title = *photo.Description
				}
				var description string
				if photo.AltDescription != nil {
					description = *photo.AltDescription
				}
				document := types.Image{
					Origin:       origin,
					OriginID:     originId,
					Extension:    extension,
					Path:         fileName,
					Width:        width,
					Height:       height,
					Title:        title,
					Description:  description,
					License:      "Unsplash License, Commercial and non-commercial purposes",
					CreationDate: &now,
					Tags:         tags,
				}

				insertedId, err := mongodb.InsertImage(collectionUnsplash, document)
				if err != nil {
					return nil, err
				}
				insertedIds = append(insertedIds, insertedId)
			}
		}

	}
	return insertedIds, nil
}

func searchPhotosPerPageUnsplash(tag string, page int) (*unsplash.PhotoSearchResult, error) {
	r := &Request{
		Host: "https://api.unsplash.com/search/photos/?",
		Args: map[string]string{
			"client_id": utils.DotEnvVariable("UNSPLASH_PUBLIC_KEY"),
			"page":      fmt.Sprint(page),
			"query":     tag,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var searchPerPage unsplash.PhotoSearchResult
	err = json.Unmarshal(body, &searchPerPage)
	if err != nil {
		return nil, err
	}
	return &searchPerPage, nil
}
