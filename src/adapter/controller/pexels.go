package controller

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"golang.org/x/exp/slices"

	"regexp"
	"strings"

	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceAdapter "scraper-backend/src/adapter/interface"
	interfaceHost "scraper-backend/src/driver/interface/host"
	model "scraper-backend/src/driver/model"
)

type ControllerPexels struct {
	Api               interfaceHost.DriverApiPexels
	ControllerPicture interfaceAdapter.ControllerPicture
	ControllerTag     interfaceAdapter.ControllerTag
	ControllerUser    interfaceAdapter.ControllerUser
}

//TODO: use model types

func (c *ControllerPexels) SearchPhotos(ctx context.Context, quality string) error {
	qualitiesAvailable := []string{"large2x", "large", "medium", "small", "portrait", "landscape", "tiny"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return fmt.Errorf("quality needs to be `large2x`(h=650), `large`(h=650), `medium`(h=350), `small`(h=130), `portrait`(h=1200), `landscape`(h=627)or `tiny`(h=200) and your is `%s`", quality)
	}

	origin := "pexels"

	searchedTags, err := c.ControllerTag.ReadTags(ctx, "searched")
	if err != nil {
		return err
	}
	// no blocked tag for pexels

	for _, searchedTag := range searchedTags {
		page := 1
		searchPerPage, err := c.Api.SearchPhotosPerPage(searchedTag.Name, page)
		if err != nil {
			return fmt.Errorf("SearchPhotosPerPage has failed: %v", err)
		}

		for page := page; page <= searchPerPage.TotalResults/searchPerPage.PerPage; page++ {
			searchPerPage, err = c.Api.SearchPhotosPerPage(searchedTag.Name, page)
			if err != nil {
				return fmt.Errorf("SearchPhotosPerPage has failed: %v", err)
			}

			for _, photo := range searchPerPage.Photos {
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
				users, err := c.ControllerUser.ReadUsers(ctx)
				if err != nil {
					return err
				}
				for _, user := range users {
					if user.OriginID == fmt.Sprint(photo.PhotographerID) {
						continue // skip the image with unwanted user
					}
				}

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

				// get buffer of image
				buffer, err := c.Api.GetFile(link)
				if err != nil {
					return fmt.Errorf("GetFile has failed: %v", err)
				}

				// image creation
				now := time.Now()
				imageSizeID := model.NewUUID()
				user := controllerModel.User{
					ID:           model.NewUUID(),
					Origin:       origin,
					Name:         photo.Photographer,
					OriginID:     fmt.Sprint(photo.PhotographerID),
					CreationDate: now,
				}
				linkURL, err := url.Parse(link)
				if err != nil {
					return err
				}
				width, err := strconv.Atoi(linkURL.Query().Get("w"))
				if err != nil {
					return err
				}
				height, err := strconv.Atoi(linkURL.Query().Get("h"))
				if err != nil {
					return err
				}
				zero := 0
				box := controllerModel.Box{
					Tlx:    zero, // original x anchor
					Tly:    zero, // original y anchor
					Width:  width,
					Height: height,
				}
				tags := map[model.UUID]controllerModel.PictureTag{
					model.NewUUID(): {
						Name:         strings.ToLower(searchedTag.Name),
						CreationDate: now,
						OriginName:   origin,
						BoxInformation: model.NewNullable(controllerModel.BoxInformation{
							ImageSizeID: imageSizeID,
							Box:         box,
						}),
					},
				}
				sizes := map[model.UUID]controllerModel.PictureSize{
					imageSizeID: {
						CreationDate: now,
						Box:          box,
					},
				}
				document := controllerModel.Picture{
					ID:           model.NewUUID(),
					Origin:       origin,
					OriginID:     fmt.Sprint(photo.ID),
					User:         user,
					Extension:    extension,
					Name:         fmt.Sprint(photo.ID),
					Sizes:        sizes,
					Title:        "",
					Description:  photo.Alt,
					License:      "No known copyright restrictions",
					CreationDate: now,
					Tags:         tags,
				}

				if c.ControllerPicture.CreatePicture(ctx, model.NewUUID(), document, buffer); err != nil {
					return fmt.Errorf("CreatePicture has failed: %v", err)
				}
			}
		}
	}
	return nil
}
