package controller

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/hbagdi/go-unsplash/unsplash"

	"strings"

	"strconv"

	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceAdapter "scraper-backend/src/adapter/interface"
	interfaceHost "scraper-backend/src/driver/interface/host"
	"scraper-backend/src/util"
	utilModel "scraper-backend/src/util/model"
)

type ControllerUnsplash struct {
	Api               interfaceHost.DriverApiUnsplash
	ControllerPicture interfaceAdapter.ControllerPicture
	ControllerTag     interfaceAdapter.ControllerTag
	ControllerUser    interfaceAdapter.ControllerUser
}

//TODO: use model types

func (c *ControllerUnsplash) SearchPhotos(ctx context.Context, quality string) error {
	qualitiesAvailable := []string{"raw", "full", "regular", "small", "thumb"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return fmt.Errorf("quality needs to be `raw`, `full`(hd), `regular`(w = 1080), `small`(w = 400) or `thumb`(w = 200) and your is `%s`", quality)
	}

	// If path is already a directory, MkdirAll does nothing and returns nil
	origin := "unsplash"

	searchedTags, err := c.ControllerTag.ReadTags(ctx, "searched")
	if err != nil {
		return err
	}

	blockedTags, err := c.ControllerTag.ReadTags(ctx, "blocked")
	if err != nil {
		return err
	}

	for _, searchedTag := range searchedTags {
		page := 1

		searchPerPage, err := c.Api.SearchPhotosPerPage(searchedTag.Name, page)
		if err != nil {
			return fmt.Errorf("searchPhotosPerPageUnsplash has failed: %v", err)
		}

		for page := page; page <= int(*searchPerPage.TotalPages); page++ {
			searchPerPage, err = c.Api.SearchPhotosPerPage(searchedTag.Name, page)
			if err != nil {
				return fmt.Errorf("searchPhotosPerPageUnsplash has failed: %v", err)
			}

			for _, photo := range *searchPerPage.Results {
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
						return err
					}
					if len(pictures) > 0 {
						continue // skip existing image
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
					return err
				}
				for _, user := range users {
					if user.OriginID == UserID {
						continue // skip the image with unwanted user
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
					continue // skip image with unwanted tag
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
				if extension == "jpeg" {
					extension = "jpg"
				}

				// get buffer of image
				buffer, err := c.Api.GetFile(link.String())
				if err != nil {
					return fmt.Errorf("GetFile has failed: %v", err)
				}

				// tags creation
				width, err := strconv.Atoi(link.Query().Get("w"))
				if err != nil {
					return err
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

				tags := make(map[uuid.UUID]controllerModel.PictureTag)
				now := time.Now()
				imageSizeID := uuid.New()
				for _, photoTag := range *photo.Tags {
					var tagTitle string
					if photoTag.Title != nil {
						tagTitle = *photoTag.Title
					}
					tags[uuid.New()] = controllerModel.PictureTag{
						Name:         strings.ToLower(tagTitle),
						CreationDate: now,
						OriginName:   origin,
						BoxInformation: utilModel.Nullable[controllerModel.BoxInformation]{
							Valid: true,
							Body: controllerModel.BoxInformation{
								ImageSizeID: imageSizeID,
								Box:         box,
							},
						},
					}
				}

				// image creation
				user := controllerModel.User{
					Origin:       origin,
					Name:         userName,
					OriginID:     UserID,
					CreationDate: now,
				}
				sizes := map[uuid.UUID]controllerModel.PictureSize{
					imageSizeID: {
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
				document := controllerModel.Picture{
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

				if c.ControllerPicture.CreatePicture(ctx, uuid.New(), document, buffer); err != nil {
					return fmt.Errorf("CreatePicture has failed: %v", err)
				}
			}
		}

	}
	return nil
}
