package gin

import (
	"context"
	serverModel "scraper-backend/src/driver/server/model"

	"github.com/google/uuid"
)

func (d DriverServerGin) CreateUserBlocked(ctx context.Context, user serverModel.User) (string, error) {
	err := d.ControllerUser.CreateUser(ctx, user.DriverUnmarshal())
	if err != nil {
		return "error", err
	}
	return "ok", nil
}

type ParamsDeleteUser struct {
	Origin string `uri:"origin" binding:"required"`
	ID     string `uri:"id" binding:"required"`
}

func (d DriverServerGin) DeleteUserBlocked(ctx context.Context, params ParamsDeleteUser) (string, error) {
	id, err := uuid.Parse(params.ID)
	if err != nil {
		return "error", err
	}
	if err := d.ControllerUser.DeleteUser(ctx, params.Origin, id); err != nil {
		return "error", err
	}
	return "ok", nil
}

func (d DriverServerGin) ReadUsers(ctx context.Context) ([]serverModel.User, error) {
	controllerUsers, err := d.ControllerUser.ReadUsers(ctx)
	if err != nil {
		return nil, err
	}
	serverUsers := make([]serverModel.User, len(controllerUsers))
	for i, controllerUser := range controllerUsers {
		serverUsers[i].DriverMarshal(controllerUser)
	}
	return serverUsers, nil
}
