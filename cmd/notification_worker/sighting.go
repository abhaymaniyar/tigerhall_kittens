package notification_worker

import (
	"context"
	"github.com/google/uuid"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
)

type sightingEmailNotifier struct {
	sightingRepo repository.SightingRepo
}

type SightingEmailNotifer interface {
	ReportSightingToAllUsers(ctx context.Context, tigerID uint, sightings []model.Sighting) error
}

func NewSightingEmailNotifer() SightingEmailNotifer {
	return &sightingEmailNotifier{
		sightingRepo: repository.NewSightingRepo(),
	}
}

func (e *sightingEmailNotifier) Process(userID uuid.UUID, data interface{}) {
	notification := Notification{
		Subject: EmailNotificationSubjectTigerSighting,
		UserID:  userID,
		Data:    data,
	}

	notificationQueue <- notification
}

type TigerSightingEmail struct {
	UserID  uint
	TigerID uint
}

// ReportSightingToAllUsers simulates reporting a tiger sighting and sends a notification email
func (e *sightingEmailNotifier) ReportSightingToAllUsers(ctx context.Context, tigerID uint, sightings []model.Sighting) error {
	reportedByUser := ctx.Value("userID")

	tigerSightingEmail := TigerSightingEmail{
		TigerID: tigerID,
	}

	for _, sighting := range sightings {
		// skipping notifications to the user who reported the tiger
		if sighting.ReportedByUserID == reportedByUser {
			continue
		}

		wg.Add(1)
		sighting := sighting
		go func() {
			defer wg.Done()
			e.Process(sighting.ReportedByUserID, tigerSightingEmail)
		}()
	}

	return nil
}
