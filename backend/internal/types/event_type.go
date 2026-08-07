package types

type EventType string

const (
	EventTypeUserLogin              EventType = "user.logged_in"
	EventTypePasswordResetRequested EventType = "user.password_reset_requested"
	EventTypePasswordChanged        EventType = "user.password_changed"
)
