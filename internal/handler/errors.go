package handler

import (
	"errors"
	"fmt"
	"tigerhall_kittens/internal/service"
	"tigerhall_kittens/internal/web"
)

func errorResponse(err error) web.ErrorInterface {
	if errors.Is(err, service.ErrTigerDoesNotExist) {
		return web.ErrBadRequest(err.Error())
	}

	if errors.Is(err, service.ErrFetchingTigerDetails) {
		return web.ErrInternalServerError(fmt.Sprintf("error while reporting sighting : %s", err.Error()))
	}

	if errors.Is(err, service.ErrFetchingExistingSightings) {
		return web.ErrInternalServerError(fmt.Sprintf("error while reporting sighting : %s", err.Error()))
	}

	if errors.Is(err, service.ErrSightingAlreadyReported) {
		return web.ErrBadRequest(fmt.Sprintf("error while reporting sighting : %s", err.Error()))
	}

	if errors.Is(err, service.ErrSendingEmailNotification) {
		return web.ErrInternalServerError(fmt.Sprintf("error while reporting sighting : %s", err.Error()))
	}

	if errors.Is(err, service.ErrInvalidUsernamePassword) {
		return web.ErrUnauthorizedRequest(fmt.Sprintf("login failed : %s", err.Error()))
	}

	if errors.Is(err, service.ErrTokenGenerationFailed) {
		return web.ErrInternalServerError(fmt.Sprintf("login failed : %s", err.Error()))
	}

	if errors.Is(err, service.ErrUserAlreadyExistsWithSameEmailUsername) {
		return web.ErrBadRequest(fmt.Sprintf("user creation failed : %s", err.Error()))
	}

	if errors.Is(err, service.ErrCreatingUser) {
		return web.ErrInternalServerError(err.Error())
	}

	return web.ErrInternalServerError(fmt.Sprintf("error while processing request : %s", err))
}
