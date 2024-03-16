package shared

type ContextKey string

const (
	CtxValueRequestId ContextKey = "request_id"
	CtxPathURL        ContextKey = "path_URL"
)
