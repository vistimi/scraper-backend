package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	typeUnsplash "github.com/hbagdi/go-unsplash/unsplash"

	"strings"

	"strconv"

	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceAdapter "scraper-backend/src/adapter/interface"
	interfaceHost "scraper-backend/src/driver/interface/host"
	model "scraper-backend/src/driver/model"
	"scraper-backend/src/util"
)

type ControllerUnsplash struct {
	Api               interfaceHost.DriverApiUnsplash
	ControllerPicture interfaceAdapter.ControllerPicture
	ControllerTag     interfaceAdapter.ControllerTag
	ControllerUser    interfaceAdapter.ControllerUser
}

//TODO: use model types

type OutputImage struct {
	OriginID *string
	Error    error
}

type InputImage struct {
	Photo       typeUnsplash.Photo
	Origin      string
	Quality     string
	BlockedTags []controllerModel.Tag
}

type OutputPage struct {
	OriginIDs []string
	Error     error
}

type InputPage struct {
	Origin      string
	Quality     string
	BlockedTags []controllerModel.Tag
	Pictures    []typeUnsplash.Photo
}

func (c *ControllerUnsplash) SearchPhotos(ctx context.Context, quality string, imageStart, imageEnd int) ([]string, error) {
	fmt.Printf("Unsplash seach images from %v to = %v\n", imageStart, imageEnd)

	qualitiesAvailable := []string{"raw", "full", "regular", "small", "thumb"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return nil, fmt.Errorf("quality needs to be `raw`, `full`(hd), `regular`(w = 1080), `small`(w = 400) or `thumb`(w = 200) and your is `%s`", quality)
	}

	origin := "unsplash"

	searchedTags, err := c.ControllerTag.ReadTags(ctx, "searched")
	if err != nil {
		return nil, err
	}
	if len(searchedTags) == 0 {
		return nil, fmt.Errorf("no searched tags")
	}

	if len(searchedTags) > 1 {
		return nil, fmt.Errorf("only one searched tag is allowed")
	}

	blockedTags, err := c.ControllerTag.ReadTags(ctx, "blocked")
	if err != nil {
		return nil, err
	}

	perPage := c.Api.GetPerPage()
	searchPerPage, err := c.Api.SearchPhotosPerPage(searchedTags[0].Name, 0)
	if err != nil {
		return nil, fmt.Errorf("searchPhotosPerPageUnsplash has failed: %v", err)
	}
	if searchPerPage.Results == nil || searchPerPage.TotalPages == nil {
		return nil, fmt.Errorf("no results found for unsplash")
	}
	total := *searchPerPage.TotalPages * perPage
	if imageStart > total || imageEnd > total {
		return nil, fmt.Errorf("indexes out of bound, total = %v", total)
	}
	fmt.Printf("last page = %v, last index = %v, per page = %v\n", *searchPerPage.TotalPages, total, perPage)

	inputsPage := make(chan InputPage)
	outputsPage := make(chan OutputPage)

	pageFrom := imageStart / perPage
	pageTo := imageEnd / perPage
	if imageEnd%perPage == 0 {
		pageTo--
	}
	fmt.Printf("page from = %v to = %v, per page = %v\n", pageFrom, pageTo, perPage)
	var wgPage sync.WaitGroup
	wgPage.Add(pageTo - pageFrom + 1)

	for page := pageFrom; page <= pageTo; page++ {
		go func() {
			for inputPage := range inputsPage {
				c.fetchPage(ctx, inputPage, outputsPage, &wgPage)
			}
		}()
	}

	// Send the inputs to the worker goroutines
	for page := pageFrom; page <= pageTo; page++ {
		searchPerPage, err := c.Api.SearchPhotosPerPage(searchedTags[0].Name, page)
		if err != nil {
			return nil, fmt.Errorf("searchPhotosPerPageUnsplash has failed: %v", err)
		}
		if searchPerPage.Results == nil {
			return nil, fmt.Errorf("no results found for unsplash for page %v", page)
		}
		pictures := *searchPerPage.Results
		ln := len(pictures)
		if page == pageFrom {
			pictures = pictures[imageStart%ln:]
		}
		if page == pageTo && imageEnd%ln != 0 {
			pictures = pictures[:imageEnd%ln]
		}
		inputsPage <- InputPage{
			Origin:      origin,
			Quality:     quality,
			BlockedTags: blockedTags,
			Pictures:    pictures,
		}
		fmt.Printf("page %v, with %v pictures\n", page, len(pictures))
	}
	close(inputsPage)

	OriginIDs := []string{}
	// Read the outputs from the output channel
	for page := pageFrom; page <= pageTo; page++ {
		outputPage := <-outputsPage
		if outputPage.Error != nil {
			return nil, outputPage.Error
		}
		OriginIDs = append(OriginIDs, outputPage.OriginIDs...)
	}
	// wgPage.Wait() // wait for all pages to be fetched

	return OriginIDs, nil
}

func (c *ControllerUnsplash) fetchPage(ctx context.Context, inputPage InputPage, outputPage chan OutputPage, wgPage *sync.WaitGroup) {
	pictures := inputPage.Pictures
	origin := inputPage.Origin
	quality := inputPage.Quality
	blockedTags := inputPage.BlockedTags

	// Init waitgroup variables
	var wgImage sync.WaitGroup // synchronize all channels
	wgImage.Add(len(pictures))

	// Set up the input and output channels
	inputsImage := make(chan InputImage)
	outputsImage := make(chan OutputImage)

	// imageStart the worker goroutines
	for i := 0; i < len(pictures); i++ {
		go func() {
			for inputImage := range inputsImage {
				c.fetchImage(ctx, inputImage, outputsImage, &wgImage)
			}
		}()
	}
	// Send the inputs to the worker goroutines
	for _, photo := range pictures {
		inputsImage <- InputImage{
			Photo:       photo,
			Origin:      origin,
			Quality:     quality,
			BlockedTags: blockedTags,
		}
	}
	close(inputsImage)

	outputPageTemp := new(OutputPage)
	// Read the results from the output channel
	for i := 0; i < len(pictures); i++ {
		outputImage := <-outputsImage

		if outputImage.Error != nil {
			outputPage <- OutputPage{OriginIDs: nil, Error: outputImage.Error}
			wgPage.Done()
			return
		}
		if outputImage.OriginID != nil {
			outputPageTemp.OriginIDs = append(outputPageTemp.OriginIDs, *outputImage.OriginID)
		}
	}
	// wgImage.Wait() // wait for all images to be fetched

	outputPage <- *outputPageTemp
	wgPage.Done()
}

func (c *ControllerUnsplash) fetchImage(ctx context.Context, inputImage InputImage, outputImage chan OutputImage, wgImage *sync.WaitGroup) {
	photo := inputImage.Photo
	origin := inputImage.Origin
	quality := inputImage.Quality
	blockedTags := inputImage.BlockedTags

	// look for existing image
	var originID string
	if photo.ID != nil {
		originID = *photo.ID
	}

	// look for existing images
	for _, state := range []string{"production", "validation", "process", "blocked"} {
		projEx := expression.NamesList(expression.Name("OriginID"))
		filtEx := expression.Name("OriginID").Equal(expression.Value(photo.ID))
		pictures, err := c.ControllerPicture.ReadPictures(ctx, state, &projEx, &filtEx)
		if err != nil {
			outputImage <- OutputImage{
				OriginID: &originID,
				Error:    err,
			}
			wgImage.Done()
			return
		}
		if len(pictures) > 0 {
			outputImage <- OutputImage{
				OriginID: nil,
				Error:    nil,
			}
			wgImage.Done()
			return // skip existing image
		}
	}

	// look for unwanted user
	var UserID string
	if photo.Photographer.ID != nil {
		UserID = *photo.Photographer.ID
	}
	var userName string
	if photo.Photographer.Username != nil {
		userName = *photo.Photographer.Username
	}
	users, err := c.ControllerUser.ReadUsers(ctx)
	if err != nil {
		outputImage <- OutputImage{
			OriginID: &originID,
			Error:    err,
		}
		wgImage.Done()
		return
	}
	for _, user := range users {
		if user.OriginID == UserID {
			outputImage <- OutputImage{
				OriginID: nil,
				Error:    nil,
			}
			wgImage.Done()
			return // skip the image with unwanted user
		}
	}

	// look for unwanted tag
	var photoTags []string
	for _, tag := range *photo.Tags {
		photoTags = append(photoTags, strings.ToLower(*tag.Title))
	}
	var blockedTagsString []string
	for _, tag := range blockedTags {
		blockedTagsString = append(blockedTagsString, tag.Name)
	}
	idx := util.FindIndexRegExp(blockedTagsString, photoTags)
	if idx != -1 {
		outputImage <- OutputImage{
			OriginID: nil,
			Error:    nil,
		}
		wgImage.Done()
		return // skip image with unwanted tag
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

	// get buffer of image
	buffer, err := c.Api.GetFile(link.String())
	if err != nil {
		outputImage <- OutputImage{
			OriginID: &originID,
			Error:    fmt.Errorf("GetFile has failed: %v", err),
		}
		wgImage.Done()
		return
	}

	// tags creation
	width, err := strconv.Atoi(link.Query().Get("w"))
	if err != nil {
		outputImage <- OutputImage{
			OriginID: &originID,
			Error:    err,
		}
		wgImage.Done()
		return
	}
	var height int
	if photo.Height != nil && photo.Width != nil {
		height = *photo.Height * width / *photo.Width
	}
	zero := 0
	box := controllerModel.Box{
		Tlx:    zero, // original x anchor
		Tly:    zero, // original y anchor
		Width:  width,
		Height: height,
	}

	tags := make([]controllerModel.PictureTag, 0, len(*photo.Tags))
	now := time.Now()
	pictureSizeID := model.NewUUID()
	for _, photoTag := range *photo.Tags {
		var tagTitle string
		if photoTag.Title != nil {
			tagTitle = *photoTag.Title
		}
		tags = append(tags, controllerModel.PictureTag{
			ID:           model.NewUUID(),
			Name:         strings.ToLower(tagTitle),
			CreationDate: now,
			OriginName:   origin,
			// BoxInformation
		})
	}

	// image creation
	user := controllerModel.User{
		ID:           model.NewUUID(),
		Origin:       origin,
		Name:         userName,
		OriginID:     UserID,
		CreationDate: now,
	}
	sizes := []controllerModel.PictureSize{
		{
			ID:           pictureSizeID,
			CreationDate: now,
			Box:          box,
		},
	}
	var title string
	if photo.Description != nil {
		title = *photo.Description
	}
	var description string
	if photo.AltDescription != nil {
		description = *photo.AltDescription
	}
	picture := controllerModel.Picture{
		Origin:       origin,
		OriginID:     originID,
		User:         user,
		Extension:    extension,
		Name:         originID,
		Sizes:        sizes,
		Title:        title,
		Description:  description,
		License:      "No known copyright restrictions",
		CreationDate: now,
		Tags:         tags,
	}

	if err := c.ControllerPicture.CreatePicture(ctx, model.NewUUID(), picture, buffer); err != nil {
		outputImage <- OutputImage{
			OriginID: &originID,
			Error:    fmt.Errorf("CreatePicture has failed: %v", err),
		}
		wgImage.Done()
		return
	}
	outputImage <- OutputImage{
		OriginID: &originID,
		Error:    nil,
	}
	wgImage.Done()
}
