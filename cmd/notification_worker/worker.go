package notification_worker

import (
	"fmt"
	"sync"
	"tigerhall_kittens/internal/logger"
	"time"
)

var (
	notificationQueue = make(chan Notification)
	wg                sync.WaitGroup
)

// StartNotificationWorker starts a goroutine to process notification emails
func StartNotificationWorker() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for notification := range notificationQueue {
			logger.D(nil, fmt.Sprintf("Sending %s notification to user %s", notification.Subject, notification.UserID))
			time.Sleep(1 * time.Second) // simulated notification sending time
		}
	}()
}

func SetupNotificationWorker() {
	defer close(notificationQueue)

	StartNotificationWorker()
	wg.Wait()
}
