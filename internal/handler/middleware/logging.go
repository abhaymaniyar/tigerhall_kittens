package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/shared"
	"time"
)

func LoggingMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		requestID := uuid.New()
		ctx := context.WithValue(r.Context(), shared.CtxValueRequestId, requestID)

		start := time.Now()
		logger.I(ctx, fmt.Sprintf("Started %s %s", r.Method, r.URL.Path), logger.Field("request_id", requestID))

		next(w, r.WithContext(ctx), ps)

		logger.I(ctx, fmt.Sprintf("Completed in %v", time.Since(start)), logger.Field("request_id", requestID))
	}
}
