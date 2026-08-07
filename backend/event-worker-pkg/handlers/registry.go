package handlers

import (
	"context"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/types"
)

type EventHandler interface {
	Handle(ctx context.Context, event *model.Event) error
}

var HandlerRegistry map[types.EventType]EventHandler

func InitiateHandlerRegistry() {
	HandlerRegistry = make(map[types.EventType]EventHandler)
	HandlerRegistry[types.EventTypeUserLogin] = UserLoggedInEventHandler()
	HandlerRegistry[types.EventTypePasswordResetRequested] = PasswordResetRequestedEventHandler()
	HandlerRegistry[types.EventTypePasswordChanged] = PasswordChangedEventHandler()
}
