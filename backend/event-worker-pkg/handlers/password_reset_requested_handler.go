package handlers

import (
	"context"

	"alwis.dev/selectify/event-worker-pkg/app"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/service"
)

type passwordResetRequestedEventHandler struct {
	userService service.UserService
}

func PasswordResetRequestedEventHandler() EventHandler {
	return &passwordResetRequestedEventHandler{
		userService: app.Service().UserService,
	}
}

func (h *passwordResetRequestedEventHandler) Handle(ctx context.Context, event *model.Event) error {
	return h.userService.ProcessPasswordResetRequested(ctx, event)
}
