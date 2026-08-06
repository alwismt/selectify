package handlers

import (
	"context"

	"alwis.dev/selectify/event-worker-pkg/app"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/service"
)

type userLoggedInEventHandler struct {
	userService service.UserService
}

func UserLoggedInEventHandler() EventHandler {
	return &userLoggedInEventHandler{
		userService: app.Service().UserService,
	}
}

func (h *userLoggedInEventHandler) Handle(ctx context.Context, event *model.Event) error {
	return h.userService.ProcessUserLoggedIn(ctx, event)
}
