package router

import (
	"bytes"
	"fmt"
	"scraper/src/mongodb"
	"scraper/src/types"
	"scraper/src/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	typeUnsplash "github.com/hbagdi/go-unsplash/unsplash"

	"encoding/json"

	"path/filepath"
	"strings"

	"strconv"
)

type ParamsSearchPhotoUnsplash struct {
	Quality    string `uri:"quality" binding:"required"`
	ImageStart int    `uri:"image_start"`
	ImageEnd   int    `uri:"image_end" binding:"required"`
}

type OutputPage struct {
	ObjectID primitive.ObjectID
	Error    error
}

type InputPage struct {
	Photo                    typeUnsplash.Photo
	Origin                   string
	Quality                  string
	UnwantedTags             []string
	S3Client                 *s3.Client
	CollectionImagesPending  *mongo.Collection
	CollectionImagesWanted   *mongo.Collection
	CollectionImagesUnwanted *mongo.Collection
	CollectionUsersUnwanted  *mongo.Collection
}

// func New(Quality string, ImageStart int, ImageEnd int) ParamsSearchPhotoUnsplash {
// 	return ParamsSearchPhotoUnsplash{"Nonel", 10, 100} // enforce the default value here
// }

func SearchPhotosUnsplash(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsSearchPhotoUnsplash) ([]primitive.ObjectID, error) {
	quality := params.Quality
	imageStart := params.ImageStart
	imageEnd := params.ImageEnd
	fmt.Printf("%+v", params)

	fmt.Printf("imageStart = %v and imageEnd = %v", imageStart, imageEnd)

	qualitiesAvailable := []string{"raw", "full", "regular", "small", "thumb"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return nil, fmt.Errorf("quality needs to be `raw`, `full`(hd), `regular`(w = 1080), `small`(w = 400) or `thumb`(w = 200) and your is `%s`", quality)
	}
	var insertedIDs []primitive.ObjectID

	// If path is already a directory, MkdirAll does nothing and returns nil
	origin := "unsplash"

	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PENDING"))
	collectionImagesWanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("PRODUCTION"))
	collectionImagesUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("UNDESIRED"))
	collectionUsersUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("USERS_UNDESIRED_COLLECTION"))

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

			// Init waitgroup variables
			var wgPage sync.WaitGroup               // synchronize all channels
			wgPage.Add(len(*searchPerPage.Results))

			// Set up the input and output channels
			inputsPage := make(chan InputPage)
			outputsPage := make(chan OutputPage)

			// Start the worker goroutines
			for i := 0; i < len(*searchPerPage.Results); i++{

				go func() {
					for inputPage := range inputsPage {
						fetchImage(inputPage, outputsPage, &wgPage)
					}
				}()

			}

			// Send the inputs to the worker goroutines
			for _, photo := range *searchPerPage.Results {
				inputsPage <- InputPage{Photo: photo, Origin: origin, Quality: quality, UnwantedTags: unwantedTags, S3Client: s3Client,	CollectionImagesPending: collectionImagesPending, CollectionImagesWanted: collectionImagesWanted, CollectionImagesUnwanted: collectionImagesUnwanted, CollectionUsersUnwanted: collectionUsersUnwanted}
			}
			close(inputsPage)


			// Read the results from the output channel
			for i := 0; i < len(*searchPerPage.Results); i++{
				outputPage := <- outputsPage

				if outputPage.Error != nil {
					return nil, outputPage.Error
				}

				insertedIDs = append(insertedIDs, outputPage.ObjectID)
			}
			wgPage.Wait()

		}

	}
	return insertedIDs, nil
}

func fetchImage(inputPage InputPage, outputPage chan OutputPage, wgPage *sync.WaitGroup) {

	photo := inputPage.Photo
	origin := inputPage.Origin
	quality := inputPage.Quality
	unwantedTags := inputPage.UnwantedTags
	s3Client := inputPage.S3Client
	collectionImagesPending := inputPage.CollectionImagesPending
	collectionImagesWanted := inputPage.CollectionImagesWanted
	collectionImagesUnwanted := inputPage.CollectionImagesUnwanted
	collectionUsersUnwanted := inputPage.CollectionUsersUnwanted

	// look for existing image
	var originID string
	if photo.ID != nil {
		originID = *photo.ID
	}
	query := bson.M{"originID": originID}
	options := options.FindOne().SetProjection(bson.M{"_id": 1})
	imagePendingFound, err := mongodb.FindOne[types.Image](collectionImagesPending, query, options)
	if err != nil {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("FindOne[Image] pending existing image has failed: %v", err),
		}
		wgPage.Done()
		return
	}
	if imagePendingFound != nil { // skip existing wanted image
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    nil,
		}
		wgPage.Done()
		return
	}
	imageWantedFound, err := mongodb.FindOne[types.Image](collectionImagesWanted, query, options)
	if err != nil {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("FindOne[Image] wanted existing image has failed: %v", err),
		}
		wgPage.Done()
		return
	}
	if imageWantedFound != nil { // skip existing pending image
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    nil,
		}
		wgPage.Done()
		return
	}
	imageUnwantedFound, err := mongodb.FindOne[types.Image](collectionImagesUnwanted, query, options)
	if err != nil {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("FindOne[Image] unwanted existing image has failed: %v", err),
		}
		wgPage.Done()
		return
	}
	if imageUnwantedFound != nil { // skip image unwanted
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    nil,
		}
		wgPage.Done()
		return
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
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("FindOne[User] has failed: %v", err),
		}
		wgPage.Done()
		return
	}
	if userFound != nil { // skip the image with unwanted user
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    nil,
		}
		wgPage.Done()
		return
	}

	// extract the photo informations
	var photoTags []string
	for _, tag := range *photo.Tags {
		photoTags = append(photoTags, strings.ToLower(*tag.Title))
	}

	// skip image if one of its tag is unwanted
	idx := utils.FindIndexRegExp(unwantedTags, photoTags) // skip image with unwated tag
	if idx != -1 {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    nil,
		}
		wgPage.Done()
		return
	}

	//find download link and extension
	var link *typeUnsplash.URL
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
	if extension == "jpeg" {
		extension = "jpg"
	}

	// get the file and rename it <id>.<format>
	fileName := fmt.Sprintf("%s.%s", *photo.ID, extension)
	path := filepath.Join(origin, fileName)

	// get buffer of image
	buffer, err := utils.GetFile(link.String())
	if err != nil {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("GetFile has failed: %v", err),
		}
		wgPage.Done()
		return
	}

	_, err = utils.UploadItemS3(s3Client, bytes.NewReader(buffer), path)
	if err != nil {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("UploadItemS3 has failed: %v", err),
		}
		wgPage.Done()
		return
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
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    err,
		}
		wgPage.Done()
		return
	}
	var height int
	if photo.Height != nil && photo.Width != nil {
		height = *photo.Height * width / *photo.Width
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
		Name:         originID,
		Size:         size,
		Title:        title,
		Description:  description,
		License:      "Unsplash License",
		CreationDate: &now,
		Tags:         tags,
	}

	insertedID, err := mongodb.InsertImage(collectionImagesPending, document)
	if err != nil {
		outputPage <- OutputPage{
			ObjectID: primitive.NilObjectID,
			Error:    fmt.Errorf("InsertImage has failed: %v", err),
		}
		wgPage.Done()
		return
	}
	outputPage <- OutputPage{
		ObjectID: insertedID,
		Error:    nil,
	}
	wgPage.Done()
}

func searchPhotosPerPageUnsplash(tag string, page int) (*typeUnsplash.PhotoSearchResult, error) {
	r := &utils.Request{
		Host: "https://api.unsplash.com/search/photos/?",
		Args: map[string]string{
			"client_id": utils.GetEnvVariable("UNSPLASH_PUBLIC_KEY"),
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

	var searchPerPage typeUnsplash.PhotoSearchResult
	err = json.Unmarshal(body, &searchPerPage)
	if err != nil {
		return nil, err
	}
	return &searchPerPage, nil
}
