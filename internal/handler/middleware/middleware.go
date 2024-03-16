package middleware

type Middleware func(nextHandler Controller) Controller

func EmptyMiddleware(nextHandler Controller) Controller {
	return nextHandler
}
