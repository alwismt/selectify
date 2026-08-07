package handlers

import (
	"context"

	"alwis.dev/selectify/event-worker-pkg/app"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/service"
)

type passwordChangedEventHandler struct {
	userService service.UserService
}

func PasswordChangedEventHandler() EventHandler {
	return &passwordChangedEventHandler{
		userService: app.Service().UserService,
	}
}

func (h *passwordChangedEventHandler) Handle(ctx context.Context, event *model.Event) error {
	return h.userService.ProcessPasswordChanged(ctx, event)
}
