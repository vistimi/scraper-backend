package router

import (
	"fmt"
	"scraper/src/mongodb"
	"scraper/src/types"
	"scraper/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hbagdi/go-unsplash/unsplash"

	"encoding/json"

	"os"

	"path/filepath"
	"strings"

	"strconv"
)

type ParamsSearchPhotoUnsplash struct {
	Quality string `uri:"quality" binding:"required"`
}

func SearchPhotosUnsplash(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsSearchPhotoUnsplash) ([]primitive.ObjectID, error) {
	quality := params.Quality
	qualitiesAvailable := []string{"raw", "full", "regular", "small", "thumb"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return nil, fmt.Errorf("quality needs to be `raw`, `full`(hd), `regular`(w = 1080), `small`(w = 400) or `thumb`(w = 200) and your is `%s`", quality)
	}
	var insertedIDs []primitive.ObjectID

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	origin := "unsplash"
	err := os.MkdirAll(filepath.Join(folderDir, origin), os.ModePerm)
	if err != nil {
		return nil, err
	}

	collectionImagesPending := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_PENDING_COLLECTION"))
	collectionImagesWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_WANTED_COLLECTION"))
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	collectionUsersUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("USERS_UNWANTED_COLLECTION"))

	unwantedTags, wantedTags, err := mongodb.TagsNames(mongoClient)
	if err != nil {
		return nil, err
	}

	for _, wantedTag := range wantedTags {
		page := 1

		searchPerPage, err := searchPhotosPerPageUnsplash(wantedTag, page)
		if err != nil {
			return nil, fmt.Errorf("searchPhotosPerPageUnsplash has failed: %v", err)
		}

		for page := page; page <= int(*searchPerPage.TotalPages); page++ {
			searchPerPage, err = searchPhotosPerPageUnsplash(wantedTag, page)
			if err != nil {
				return nil, fmt.Errorf("searchPhotosPerPageUnsplash has failed: %v", err)
			}

			for _, photo := range *searchPerPage.Results {
				// look for existing image
				var originID string
				if photo.ID != nil {
					originID = *photo.ID
				}
				query := bson.M{"originID": originID}
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
				var userName string
				if photo.Photographer.Username != nil {
					userName = *photo.Photographer.Username
				}
				var UserID string
				if photo.Photographer.ID != nil {
					UserID = *photo.Photographer.ID
				}
				query = bson.M{"origin": origin,
					"$or": bson.A{
						bson.M{"originID": UserID},
						bson.M{"name": userName},
					},
				}
				userFound, err := mongodb.FindOne[types.User](collectionUsersUnwanted, query)
				if err != nil {
					return nil, fmt.Errorf("FindOne[User] has failed: %v", err)
				}
				if userFound != nil {
					continue // skip the image with unwanted user
				}

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
				var link *unsplash.URL
				switch quality {
				case "raw":
					link = photo.Urls.Raw
				case "full":
					link = photo.Urls.Full
				case "regular":
					link = photo.Urls.Regular
				case "small":
					link = photo.Urls.Small
				case "thumb":
					link = photo.Urls.Thumb
				}
				extension := link.Query().Get("fm")

				// get the file and rename it <id>.<format>
				fileName := fmt.Sprintf("%s.%s", *photo.ID, extension)
				path := filepath.Join(origin, fileName)

				_, err = UploadS3(s3Client, link.String(), path)
				if err != nil {
					return nil, fmt.Errorf("UploadS3 has failed: %v", err)
				}

				// tags creation
				var tags []types.Tag
				now := time.Now()
				imageSizeID := primitive.NewObjectID()
				tagOrigin := types.TagOrigin{
					Name:        origin,
					ImageSizeID: imageSizeID,
				}
				for _, photoTag := range *photo.Tags {
					var tagTitle string
					if photoTag.Title != nil {
						tagTitle = *photoTag.Title
					}
					tag := types.Tag{
						Name:         strings.ToLower(tagTitle),
						Origin:       tagOrigin,
						CreationDate: &now,
					}
					tags = append(tags, tag)
				}

				// image creation
				user := types.User{
					Origin:       origin,
					Name:         userName,
					OriginID:     UserID,
					CreationDate: &now,
				}
				width, err := strconv.Atoi(link.Query().Get("w"))
				if err != nil {
					return nil, err
				}
				var height int
				if photo.Height != nil && photo.Width != nil {
					height = *photo.Height * width / *photo.Width
				}
				zero := 0
				box := types.Box{
					X:      &zero, // original x anchor
					Y:      &zero, // original y anchor
					Width:  &width,
					Height: &height,
				}
				size := []types.ImageSize{{
					ID:           imageSizeID,
					CreationDate: &now,
					Box:          box,
				}}
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
					OriginID:     originID,
					User:         user,
					Extension:    extension,
					Name:         fileName,
					Size:         size,
					Title:        title,
					Description:  description,
					License:      "Unsplash License",
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

func searchPhotosPerPageUnsplash(tag string, page int) (*unsplash.PhotoSearchResult, error) {
	r := &Request{
		Host: "https://api.unsplash.com/search/photos/?",
		Args: map[string]string{
			"client_id": utils.DotEnvVariable("UNSPLASH_PUBLIC_KEY"),
			"per_page":  "80", // default 10
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
