package router

import (
	"bytes"
	"fmt"
	"regexp"
	"scraper/src/mongodb"
	"scraper/src/types"
	"scraper/src/utils"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"encoding/json"
	"net/url"
	"path/filepath"
)

type ParamsSearchPhotoPexels struct {
	Quality string `uri:"quality" binding:"required"`
}

func SearchPhotosPexels(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsSearchPhotoPexels) (interface{}, error) {
	quality := params.Quality
	qualitiesAvailable := []string{"large2x", "large", "medium", "small", "portrait", "landscape", "tiny"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return nil, fmt.Errorf("quality needs to be `large2x`(h=650), `large`(h=650), `medium`(h=350), `small`(h=130), `portrait`(h=1200), `landscape`(h=627)or `tiny`(h=200) and your is `%s`", quality)
	}
	var insertedIDs []primitive.ObjectID

	origin := "pexels"

	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))
	collectionImagesWanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PRODUCTION"))
	collectionImagesUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("UNDESIRED"))
	collectionUsersUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("USERS_UNDESIRED_COLLECTION"))

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
				query := bson.M{"originID": fmt.Sprint(photo.ID)}
				options := options.FindOne().SetProjection(bson.M{"_id": 1})
				imagePendingFound, err := mongodb.FindOne[types.Image](collectionImagesPending, query, options)
				if err != nil {
					return nil, fmt.Errorf("FindOne[Image] pending existing image has failed: %v", err)
				}
				if imagePendingFound != nil {
					continue // skip existing wanted image
				}
				imageWantedFound, err := mongodb.FindOne[types.Image](collectionImagesWanted, query, options)
				if err != nil {
					return nil, fmt.Errorf("FindOne[Image] wanted existing image has failed: %v", err)
				}
				if imageWantedFound != nil {
					continue // skip existing pending image
				}
				imageUnwantedFound, err := mongodb.FindOne[types.Image](collectionImagesUnwanted, query, options)
				if err != nil {
					return nil, fmt.Errorf("FindOne[Image] unwanted existing image has failed: %v", err)
				}
				if imageUnwantedFound != nil {
					continue // skip image unwanted
				}

				// look for unwanted Users
				query = bson.M{"origin": origin,
					"$or": bson.A{
						bson.M{"originID": fmt.Sprint(photo.PhotographerID)},
						bson.M{"name": photo.Photographer},
					},
				}
				userFound, err := mongodb.FindOne[types.User](collectionUsersUnwanted, query)
				if err != nil {
					return nil, fmt.Errorf("FindOne[User] has failed: %v", err)
				}
				if userFound != nil {
					continue // skip the image with unwanted user
				}

				// no tags present so no need to check for unwanted tags

				//find download link and extension
				var link string
				switch quality {
				case "large2x":
					link = photo.Src.Large2X
				case "large":
					link = photo.Src.Large
				case "medium":
					link = photo.Src.Medium
				case "small":
					link = photo.Src.Small
				case "portrait":
					link = photo.Src.Portrait
				case "landscape":
					link = photo.Src.Landscape
				case "tiny":
					link = photo.Src.Tiny
				}
				regexpMatch := regexp.MustCompile(`\.\w+\?`) // matches a word  preceded by `.` and followed by `?`
				extension := string(regexpMatch.Find([]byte(link)))
				extension = extension[1 : len(extension)-1] // remove the `.` and `?` because retgexp hasn't got assertions
				if extension == "jpeg" {
					extension = "jpg"
				}

				// get the file and rename it <id>.<format>
				fileName := fmt.Sprintf("%d.%s", photo.ID, extension)
				path := filepath.Join(origin, fileName)

				// get buffer of image
				buffer, err := utils.GetFile(link)
				if err != nil {
					return nil, fmt.Errorf("GetFile has failed: %v", err)
				}

				_, err = utils.UploadItemS3(s3Client, bytes.NewReader(buffer), path)
				if err != nil {
					return nil, fmt.Errorf("UploadItemS3 has failed: %v", err)
				}

				// image creation
				now := time.Now()
				imageSizeID := primitive.NewObjectID()
				tagOrigin := types.TagOrigin{
					Name:        origin,
					ImageSizeID: imageSizeID,
				}
				tags := []types.Tag{
					{
						Name:         wantedTag,
						Origin:       tagOrigin,
						CreationDate: &now,
					},
				}
				user := types.User{
					Origin:       origin,
					Name:         photo.Photographer,
					OriginID:     fmt.Sprint(photo.PhotographerID),
					CreationDate: &now,
				}
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
				zero := 0
				box := types.Box{
					Tlx:    &zero, // original x anchor
					Tly:    &zero, // original y anchor
					Width:  &width,
					Height: &height,
				}
				size := []types.ImageSize{{
					ID:           imageSizeID,
					CreationDate: &now,
					Box:          box,
				}}
				document := types.Image{
					Origin:       origin,
					OriginID:     fmt.Sprint(photo.ID),
					User:         user,
					Extension:    extension,
					Name:         fmt.Sprint(photo.ID),
					Size:         size,
					Title:        "",
					Description:  photo.Alt,
					License:      "Pexels License",
					CreationDate: &now,
					Tags:         tags,
				}
				insertedID, err := mongodb.InsertImage(collectionImagesPending, document)
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
	r := &utils.Request{
		Host: "https://api.pexels.com/v1/search?",
		Args: map[string]string{
			"query":    tag,
			"per_page": "80", // default 15, max 80
			"page":     fmt.Sprint(page),
		},
		Header: map[string][]string{
			"Authorization": {utils.GetEnvVariable("PEXELS_PUBLIC_KEY")},
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
