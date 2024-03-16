package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/shared"
	"tigerhall_kittens/internal/web"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	APIVersionV1 = 1
)

type ResponseBuilder func(data *web.JSONResponse, responseErr web.ErrorInterface) *web.JSONResponse

func ServeV1Endpoint(middleware Middleware, handler Controller) httprouter.Handle {
	return serve(buildResponseBuilder(APIVersionV1), middleware(handler))
}

func serve(responseBuilder ResponseBuilder, handler Controller) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		startTime := time.Now()
		requestId := RequestHeaderId(req)
		contextWithResult := context.WithValue(req.Context(), shared.CtxValueRequestId, requestId)

		if getURL(req) != "" {
			contextWithResult = context.WithValue(contextWithResult, shared.CtxPathURL, getURL(req))
		}

		req = req.WithContext(contextWithResult)

		webReq := web.NewRequest(req)
		for i := range ps {
			webReq.SetPathParam(ps[i].Key, ps[i].Value)
		}

		reqBody, decodeErr := readRequestBody(req)
		if decodeErr != nil {
			logger.I(req.Context(), "Error while decoding the request body",
				logger.Field("error", decodeErr.Error()))
		}

		defer func() {
			if recvr := recover(); recvr != nil {
				errorMessage := fmt.Sprintf("%v", recvr)
				err := web.ErrInternalServerError(errorMessage)
				w.WriteHeader(err.HTTPStatusCode())
				writeResponse(req.Context(), w, responseBuilder(nil, err))
				logger.E(req.Context(), err, "Request failed",
					logger.Field("error", errorMessage),
					logger.Field("status", err.HTTPStatusCode()),
					logger.Field("path", getURL(req)),
					logger.Field("request_params", req.URL.Query()),
					logger.Field("request_body", reqBody),
					logger.Field("duration_ms", float64(time.Since(startTime).Milliseconds())),
					logger.Field("method", req.Method),
					logger.Field("stack", string(debug.Stack())),
				)
			}
		}()

		data, responseErr := handler(&webReq)
		responseCode := responseCode(responseErr)

		// setting response headers
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(responseCode)

		writeResponse(req.Context(), w, responseBuilder(data, responseErr))

		if responseCode >= http.StatusInternalServerError {
			logger.E(req.Context(), responseErr, "Request failed",
				logger.Field("error_cause", responseErr.Cause()),
				logger.Field("status", responseErr.HTTPStatusCode()),
				logger.Field("path", getURL(req)),
				logger.Field("request_params", req.URL.Query()),
				logger.Field("request_body", reqBody),
				logger.Field("duration_ms", float64(time.Since(startTime).Milliseconds())),
				logger.Field("method", req.Method),
				logger.Field("stack", string(debug.Stack())),
			)

			return
		}

		if responseCode >= http.StatusBadRequest {
			logger.W(req.Context(), "Client request error",
				logger.Field("error", responseErr.Error()),
				logger.Field("status", responseErr.HTTPStatusCode()),
				logger.Field("path", getURL(req)),
				logger.Field("request_params", req.URL.Query()),
				logger.Field("request_body", reqBody),
				logger.Field("duration_ms", float64(time.Since(startTime).Milliseconds())),
				logger.Field("method", req.Method),
				logger.Field("stack", string(debug.Stack())),
			)

			return
		}

		logger.I(req.Context(), "Request processed",
			logger.Field("status", responseCode),
			logger.Field("path", getURL(req)),
			logger.Field("request_params", webReq.QueryParams()),
			logger.Field("duration_ms", float64(time.Since(startTime).Milliseconds())),
			logger.Field("method", req.Method),
			logger.Field("request_body", reqBody),
		)
	}
}

func RequestHeaderId(req *http.Request) string {
	var requestId string
	var err error

	requestId = uuid.New().String()
	if err != nil {
		logger.E(req.Context(), err, "endpoint/RequestHeaderId")
	}
	return requestId
}

func writeResponse(ctx context.Context, w http.ResponseWriter, response *web.JSONResponse) {
	_, err := w.Write(response.ByteArray(ctx))
	if err != nil {
		logger.E(ctx, err, "error in writing response", logger.Field("error", err.Error()))
	}
}

func buildResponseBuilder(version int) ResponseBuilder {
	return func(data *web.JSONResponse, err web.ErrorInterface) *web.JSONResponse {
		if err == nil {
			return successResponse(version, data)
		} else {
			return errorResponse(version, err)
		}
	}
}

func successResponse(version int, data *web.JSONResponse) *web.JSONResponse {
	return &web.JSONResponse{
		"success":     true,
		"data":        data,
		"api_version": version,
	}
}

func errorResponse(version int, err web.ErrorInterface) *web.JSONResponse {
	return &web.JSONResponse{
		"success": false,
		"error": map[string]interface{}{
			"code":    err.Code(),
			"message": err.Description(),
		},
		"api_version": version,
	}
}

func responseCode(err web.ErrorInterface) int {
	if err != nil {
		return err.HTTPStatusCode()
	}
	return http.StatusOK
}

func getURL(r *http.Request) string {
	url := "undefined"
	if r != nil && r.URL != nil && r.URL.Path != "" {
		url = r.URL.Path
	}
	return url
}

func readRequestBody(req *http.Request) (map[string]interface{}, error) {
	var jsonPayload map[string]interface{}

	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		if len(body) != 0 {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			err := json.Unmarshal(body, &jsonPayload)
			if err != nil {
				return nil, err
			}
		}
	}

	return jsonPayload, nil
}
