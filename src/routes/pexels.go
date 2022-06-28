package routes

import (
	"fmt"
	"regexp"
	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"os"

	"encoding/json"
	"net/url"
	"path/filepath"
)

func SearchPhotosPexels(mongoClient *mongo.Client) (interface{}, error) {
	var insertedIDs []primitive.ObjectID

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	origin := "pexels"
	err := os.MkdirAll(filepath.Join(folderDir, origin), os.ModePerm)
	if err != nil {
		return nil, err
	}

	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))

	_, wantedTags, err := mongodb.TagsNames(mongoClient)
	if err != nil {
		return nil, err
	}

	for _, wantedTag := range wantedTags {
		page := 1
		searchPerPage, err := searchPhotosPerPagePexels(wantedTag, page)
		if err != nil {
			return nil, fmt.Errorf("searchPhotosPerPagePexels has failed: %v", err)
		}

		for page := page; page <= searchPerPage.TotalResults/searchPerPage.PerPage; page++ {
			searchPerPage, err = searchPhotosPerPagePexels(wantedTag, page)
			if err != nil {
				return nil, fmt.Errorf("searchPhotosPerPagePexels has failed: %v", err)
			}

			for _, photo := range searchPerPage.Photos {
				// look for existing image
				_, err := mongodb.FindImageIDByOriginID(collectionImages, fmt.Sprint(photo.ID))
				if err != nil {
					return nil, fmt.Errorf("FindImageIDByOriginID has failed: %v", err)
				}

				//find download link and extension
				link := photo.Src.Large
				regexpMatch := regexp.MustCompile(`\.\w+\?`) // matches a word  preceded by `.` and followed by `?`
				extension := string(regexpMatch.Find([]byte(link)))
				extension = extension[1 : len(extension)-1] // remove the `.` and `?` because retgexp hasn't got assertions

				// download photo into folder and rename it <id>.<format>
				fileName := fmt.Sprintf("%d.%s", photo.ID, extension)
				path := fmt.Sprintf(filepath.Join(folderDir, origin, fileName))
				err = DownloadFile(link, path)
				if err != nil {
					return nil, fmt.Errorf("DownloadFile has failed: %v", err)
				}

				// tags creation
				now := time.Now()
				tags := []types.Tag{
					{
						Name:         wantedTag,
						Origin:       origin,
						CreationDate: &now,
					},
				}

				// user creation
				user := types.User{
					Origin:       origin,
					Name:         photo.Photographer,
					OriginID:     fmt.Sprint(photo.PhotographerID),
					CreationDate: &now,
				}

				// image creation
				linkURL, err := url.Parse(link)
				if err != nil {
					return nil, err
				}
				width, err := strconv.Atoi(linkURL.Query().Get("w"))
				if err != nil {
					return nil, err
				}
				height, err := strconv.Atoi(linkURL.Query().Get("h"))
				if err != nil {
					return nil, err
				}
				document := types.Image{
					Origin:       origin,
					OriginID:     fmt.Sprint(photo.ID),
					User:         user,
					Extension:    extension,
					Path:         fileName,
					Width:        width,
					Height:       height,
					Title:        "",
					Description:  photo.Alt,
					License:      "Pexels License",
					CreationDate: &now,
					Tags:         tags,
				}

				insertedID, err := mongodb.InsertImage(collectionImages, document)
				if err != nil {
					return nil, fmt.Errorf("InsertImage has failed: %v", err)
				}
				insertedIDs = append(insertedIDs, insertedID)
			}
		}
	}
	return insertedIDs, nil
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
	Alt             string       `json:"alt"`
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
			"page":  fmt.Sprint(page),
		},
		Header: map[string][]string{
			"Authorization": {utils.DotEnvVariable("PEXELS_PUBLIC_KEY")},
		},
	}
	// fmt.Println(r.URL())

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
