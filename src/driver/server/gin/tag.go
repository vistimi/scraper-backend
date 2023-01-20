package gin

import (
	"context"
	serverModel "scraper-backend/src/driver/server/model"

	"github.com/google/uuid"
)

func (d DriverServerGin) CreateTag(ctx context.Context, tag serverModel.Tag) (string, error) {
	tag.Type = "scraped"
	err := d.ControllerTag.CreateTag(ctx, tag.DriverUnmarshal())
	if err != nil {
		return "error", err
	}
	return "ok", nil
}

func (d DriverServerGin) CreateTagBlocked(ctx context.Context, tag serverModel.Tag) (string, error) {
	tag.Type = "blocked"
	err := d.ControllerTag.CreateTag(ctx, tag.DriverUnmarshal())
	if err != nil {
		return "error", err
	}
	return "ok", nil
}

type ParamsDeleteTag struct {
	ID   uuid.UUID `uri:"id" binding:"required"`
}

func (d DriverServerGin) DeleteTag(ctx context.Context, params ParamsDeleteTag) (string, error) {
	err := d.ControllerTag.DeleteTag(ctx, "scraped", params.ID)
	if err != nil {
		return "error", err
	}
	return "ok", nil
}

func (d DriverServerGin) DeleteTagBlocked(ctx context.Context, params ParamsDeleteTag) (string, error) {
	err := d.ControllerTag.DeleteTag(ctx, "blocked", params.ID)
	if err != nil {
		return "error", err
	}
	return "ok", nil
}

func (d DriverServerGin) ReadTags(ctx context.Context) ([]serverModel.Tag, error) {
	controllerTags, err := d.ControllerTag.ReadTags(ctx, "scraped")
	if err != nil {
		return nil, err
	}
	serverTags := make([]serverModel.Tag, len(controllerTags))
	for i, controllerTag := range controllerTags{
		serverTags[i].DriverMarshal(controllerTag)
	}
	return serverTags, nil
}

func (d DriverServerGin) ReadTagsBlocked(ctx context.Context) ([]serverModel.Tag, error) {
	controllerTags, err := d.ControllerTag.ReadTags(ctx, "blocked")
	if err != nil {
		return nil, err
	}
	serverTags := make([]serverModel.Tag, len(controllerTags))
	for i, controllerTag := range controllerTags{
		serverTags[i].DriverMarshal(controllerTag)
	}
	return serverTags, nil
}