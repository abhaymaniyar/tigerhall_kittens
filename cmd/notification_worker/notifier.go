package notification_worker

import "github.com/google/uuid"

const EmailNotificationSubjectTigerSighting = "Tiger Sighting Email"

type Notification struct {
	Subject string
	UserID  uuid.UUID
	Data    interface{}
}

type Notifier interface {
	Process(userID uuid.UUID, data interface{}) error
}
