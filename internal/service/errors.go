package service

import "errors"

var (
	ErrTigerDoesNotExist    = errors.New("tiger does not exist")
	ErrFetchingTigerDetails = errors.New("unable to fetch tiger details")

	ErrFetchingExistingSightings = errors.New("unable to check existing sightings")
	ErrSightingAlreadyReported   = errors.New("already reported in range")

	ErrSendingEmailNotification = errors.New("unable to send email notifications")

	ErrInvalidUsernamePassword = errors.New("invalid username or password")
	ErrTokenGenerationFailed   = errors.New("failed to generate token")

	ErrUserAlreadyExistsWithSameEmailUsername = errors.New("user already exists with same email/username")
	ErrCreatingUser                           = errors.New("error while creating user")
)
