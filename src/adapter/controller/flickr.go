package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/foolin/pagser"
	"github.com/google/uuid"

	"golang.org/x/exp/slices"

	"regexp"
	"strings"

	controllerModel "scraper-backend/src/adapter/controller/model"
	interfaceAdapter "scraper-backend/src/adapter/interface"
	driverHost "scraper-backend/src/driver/host"
	interfaceHost "scraper-backend/src/driver/interface/host"
	"scraper-backend/src/util"
	utilModel "scraper-backend/src/util/model"
	hostModel "scraper-backend/src/driver/host/model"
)

type ControllerFlickr struct {
	Api               interfaceHost.DriverApiFlickr
	ControllerPicture interfaceAdapter.ControllerPicture
	ControllerTag     interfaceAdapter.ControllerTag
	ControllerUser    interfaceAdapter.ControllerUser
}

//TODO: use model types

// Find all the photos with specific quality and folder directory.
func (c *ControllerFlickr) SearchPhotos(ctx context.Context, params driverHost.ParamsSearchPhotoFlickr) error {
	quality := params.Quality
	qualitiesAvailable := []string{"Small", "Medium", "Large", "Original"}
	idx := slices.IndexFunc(qualitiesAvailable, func(qualityAvailable string) bool { return qualityAvailable == quality })
	if idx == -1 {
		return fmt.Errorf("quality needs to be `Original`(w=2400), `Large`(w=1024), `Medium`(w = 500) or `Small`(w = 240) and your is `%s`", quality)
	}
	parser := pagser.New() // parsing html in string responses

	origin := "flickr"

	searchedTags, err := c.ControllerTag.ReadTags(ctx, "searched")
	if err != nil {
		return err
	}

	blockedTags, err := c.ControllerTag.ReadTags(ctx, "blocked")
	if err != nil {
		return err
	}

	for _, searchedTag := range searchedTags {

		// all the commercial use licenses
		// https://www.flickr.com/services/api/flickr.photos.licenses.getInfo.html
		var licenseIDsNames = map[string]string{
			"4":  "Attribution License",
			"5":  "Attribution-ShareAlike License",
			"7":  "No known copyright restrictions",
			"9":  "Public Domain Dedication (CC0)",
			"10": "Public Domain Mark",
		}
		licenseIDs := [5]string{"4", "5", "7", "9", "10"}
		for _, licenseID := range licenseIDs {

			// start with the first page
			page := 1
			searchPerPage, err := c.Api.SearchPhotosPerPage(parser, licenseID, searchedTag.Name, fmt.Sprint(page))
			if err != nil {
				return fmt.Errorf("SearchPhotosPerPage has failed: %v", err)
			}

			for page := page; page <= int(searchPerPage.Pages); page++ {
				searchPerPage, err := c.Api.SearchPhotosPerPage(parser, licenseID, searchedTag.Name, fmt.Sprint(page))
				if err != nil {
					return fmt.Errorf("searchPhotosPerPageFlickr has failed: %v", err)
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

					// extract the photo informations
					infoData, err := c.Api.InfoPhoto(parser, photo)
					if err != nil {
						return fmt.Errorf("InfoPhoto has failed: %v", err)
					}
					if infoData.OriginalFormat == "jpeg" {
						infoData.OriginalFormat = "jpg"
					}

					// look for unwanted user
					users, err := c.ControllerUser.ReadUsers(ctx)
					if err != nil {
						return err
					}
					for _, user := range users {
						if user.OriginID == infoData.UserID {
							continue // skip the image with unwanted user
						}
					}

					// look for unwanted tag
					var photoTags []string
					for _, tag := range infoData.Tags {
						photoTags = append(photoTags, strings.ToLower(tag.Name))
					}
					var blockedTagsString []string
					for _, tag := range blockedTags {
						blockedTagsString = append(blockedTagsString, tag.Name)
					}
					idx := util.FindIndexRegExp(blockedTagsString, photoTags)
					if idx != -1 {
						continue // skip image with unwanted tag
					}

					// extract the photo download link
					downloadData, err := c.Api.DownloadPhoto(parser, photo.ID)
					if err != nil {
						return fmt.Errorf("DownloadPhoto has failed: %v", err)
					}

					// get the download link for the correct resolution
					label := strings.ToLower(quality)
					regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, label)
					idx = slices.IndexFunc(downloadData.Photos, func(download hostModel.DownloadPhotoSingleData) bool { return strings.ToLower(download.Label) == label })
					if idx == -1 {
						idx = slices.IndexFunc(downloadData.Photos, func(download hostModel.DownloadPhotoSingleData) bool {
							matched, err := regexp.Match(regexpMatch, []byte(strings.ToLower(download.Label)))
							if err != nil {
								return false
							}
							return matched
						})
					}
					if idx == -1 {
						return fmt.Errorf("cannot find label %s and its derivatives %s in SearchPhoto! id %s has available the following:%v", label, regexpMatch, photo.ID, downloadData)
					}

					// get buffer of image
					buffer, err := c.Api.GetFile(downloadData.Photos[idx].Source)
					if err != nil {
						return fmt.Errorf("GetFile has failed: %v", err)
					}

					// image creation
					imageSizeID := uuid.New()
					tags := make(map[uuid.UUID]controllerModel.PictureTag)
					now := time.Now()
					user := controllerModel.User{
						Origin:       origin,
						Name:         infoData.UserName,
						OriginID:     infoData.UserID,
						CreationDate: now,
					}
					zero := 0
					box := controllerModel.Box{
						Tlx:    zero, // original x anchor
						Tly:    zero, // original y anchor
						Width:  downloadData.Photos[idx].Width,
						Height: downloadData.Photos[idx].Height,
					}
					for _, infoDataTag := range infoData.Tags {
						tags[uuid.New()] = controllerModel.PictureTag{
							Name:         strings.ToLower(infoDataTag.Name),
							CreationDate: now,
							OriginName:   origin,
							BoxInformation: utilModel.Nullable[controllerModel.BoxInformation]{
								Valid: true,
								Body: controllerModel.BoxInformation{
									ImageSizeID: imageSizeID,
									Box: box,
								},
							},
						}
					}
					sizes := map[uuid.UUID]controllerModel.PictureSize{
						imageSizeID: {
							CreationDate: now,
							Box:          box,
						},
					}
					document := controllerModel.Picture{
						Origin:       origin,
						OriginID:     photo.ID,
						User:         user,
						Extension:    infoData.OriginalFormat,
						Name:         photo.ID,
						Sizes:        sizes,
						Title:        infoData.Title,
						Description:  infoData.Description,
						License:      licenseIDsNames[licenseID],
						CreationDate: now,
						Tags:         tags,
					}

					if c.ControllerPicture.CreatePicture(ctx, uuid.New(), document, buffer); err != nil {
						return fmt.Errorf("CreatePicture has failed: %v", err)
					}
				}
			}
		}
	}
	return nil
}
