package middlewares

import (
	"gitlab.com/deroo/gocart/errors"
	"gitlab.com/deroo/gocart/helpers"
	"strings"
	"time"

	"github.com/sirupsen/logrus"


	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

// RequestContext is used to store the user request info
type RequestContext struct {
	requestID      string
	crossRequestID string
	Token          *helpers.JWTToken
	Log            *logrus.Entry
	start          time.Time
}

func GetRequestContext(c echo.Context) *RequestContext {
	return c.(*customContext).requestContext
}

type customContext struct {
	echo.Context
	requestContext *RequestContext
}

// TODO: split token actions from this middleware
// TODO: change logrus.Entry to interface
// IDEA: criar token de sessão para ter um tracking de todos os requests?

func RequestContextMiddleware(logger *logrus.Entry,
	decodeToken func(encodedToken string) (*helpers.JWTToken, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rc := &RequestContext{}

			// Start time
			// IDEA: this can be used to calculate request internal time
			rc.start = time.Now()
			uuidV4, _ := uuid.NewV4()

			// Generate request ID and get CrossRequestID
			rc.requestID = strings.Replace(uuidV4.String(), "-", "", -1)
			if c.Request().Header.Get("X-Request-ID") != "" {
				rc.crossRequestID = c.Request().Header.Get("X-Request-ID")
			}

			// Creating context logger with fields
			e := logger.WithFields(logrus.Fields{
				"requestId": rc.requestID,
			})
			if rc.crossRequestID != "" {
				e = e.WithField("crossRequestId", rc.crossRequestID)
			}
			rc.Log = e

			// Parse token
			var authToken string
			if c.Request().Header.Get("Authorization") == "" && c.QueryParam("authorization") == "" {
				return errors.NewUnauthorizedError()
			}
			if c.Request().Header.Get("Authorization") != "" {
				stkn := strings.Split(c.Request().Header.Get("Authorization"), " ")
				if len(stkn) != 2 {
					err := errors.NewUnauthorizedError()
					err.DeveloperMessage = "Authorization type not found. Maybe you forgot to include 'Bearer' as a prefix for the token"
					return err
				}
				authToken = strings.Split(c.Request().Header.Get("Authorization"), " ")[1]
			}
			if c.QueryParam("authorization") != "" {
				authToken = c.QueryParam("authorization")
			}
			jwtToken, err := decodeToken(authToken)
			if err != nil {
				return errors.NewUnauthorizedError()
			}
			rc.Token = jwtToken

			cc := &customContext{c, rc}
			return next(cc)
		}
	}
}
