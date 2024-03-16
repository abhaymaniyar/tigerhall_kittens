package middleware

import (
	"tigerhall_kittens/internal/web"
)

type Controller func(request *web.Request) (*web.JSONResponse, web.ErrorInterface)
