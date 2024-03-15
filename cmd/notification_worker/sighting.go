package notification_worker

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/repository"
)

type sightingEmailNotifier struct {
	sightingRepo repository.SightingRepo
}

type SightingEmailNotifer interface {
	ReportSightingToAllUsers(ctx context.Context, tigerID uint) error
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
func (e *sightingEmailNotifier) ReportSightingToAllUsers(ctx context.Context, tigerID uint) error {
	reportedByUser := ctx.Value("userID")

	tigerSightingEmail := TigerSightingEmail{
		TigerID: tigerID,
	}

	sightings, err := e.sightingRepo.GetSightings(ctx, repository.GetSightingOpts{
		TigerID: tigerID,
	})

	if err != nil {
		logger.E(ctx, err, "Error while fetching existing sightings for the tiger", logger.Field("tiger_id", tigerID))
		return errors.New("error while fetching existing sightings")
	}

	for _, sighting := range *sightings {
		// skipping notifications to the user who reported the tiger
		if sighting.ReportedByUserID.String() == reportedByUser.(string) {
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
