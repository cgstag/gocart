package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rifflock/lfshook"
	"github.com/shopspring/decimal"

	"./config"
	"./middlewares"
	"./stores/datastores"

	"github.com/Sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"./errors"
	"./handlers"
	"./helpers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Version specifies the current version of the application.
// The value of this variable is replaced with the latest git tag
// by "make" while building the API
var Version string

// config stores the application-wide configurations
func main() {
	// load application configurations
	config := config.MustLoadConfig()

	// creating logger
	logger := initLogrus(config)

	// set decimal configs
	decimal.DivisionPrecision = 2
	decimal.MarshalJSONWithoutQuotes = true

	// DB Connection
	if err := datastores.SetupDatastore(string(config.DB.Dialect),
		config.DB.DSN,
		config.DB.MaxIdleConnections,
		config.DB.MaxOpenConnections,
		config.Debug,
	); err != nil {
		panic(fmt.Errorf("api: error when loading database connection. %s", err))
	}
	defer datastores.Database.Close()

	// Setup cache
	err := handlers.SetupCache(config.Cache.DefaultExpire, config.Cache.DefaultPurge)
	if err != nil {
		panic(fmt.Errorf("api: error when loading Cache : %s", err))
	}

	// Init router
	router := initRouter(config, logger)

	// start server
	address := fmt.Sprintf("%v:%v", config.Host, config.Port)
	logger.Infof("Server version [%s] is started at %v\n", Version, address)

	panic(router.Start(address))
}

func initLogrus(config *config.ApiConfig) *logrus.Entry {
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{}

	// Disable log to stdout
	if !config.Log.Outputs.StdoutEnable {
		l.Out = ioutil.Discard
	}

	// Log to file
	if config.Log.Outputs.File.Enabled {
		var llmap lfshook.PathMap
		llmap = make(map[logrus.Level]string)
		for _, level := range logrus.AllLevels {
			llmap[level] = config.Log.Outputs.File.Path
		}
		l.Hooks.Add(lfshook.NewHook(llmap))
	}

	// Log level
	lvl, err := logrus.ParseLevel(config.Log.Level)
	if err != nil {
		l.WithError(err).Error("api: invalid log level... fallback to debug level")
		lvl = logrus.DebugLevel
	}
	l.Level = lvl
	return l.WithField("version", Version)
}

func initRouter(config *config.ApiConfig, logger *logrus.Entry) *echo.Echo {
	router := echo.New()
	// When running in debug mode,the returned JSON is always "pretty printed"
	router.Debug = config.Debug

	helpers.SetupTokenHelper(config.SessionToken.Secret,
		config.SessionToken.Issuer,
		360, // Store Token lasts 1 year
		config.SessionToken.Duration,
		config.SessionToken.Duration,
	)

	// Healthcheck
	router.Match([]string{"GET", "HEAD"}, "/health", func(c echo.Context) error {
		return c.String(http.StatusOK, Version)
	})

	// Set echo.Logger to use primary logger
	router.Logger.SetOutput(logger.Writer())

	// Removing unnecessary trailing slash at the end of the path
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize: 1 << 10, // 1 KB
		}),
	)

	router.HTTPErrorHandler = func(err error, c echo.Context) {
		switch e := err.(type) {
		case *errors.APIError:
			l := logger.
				WithField("http_status", e.HTTPStatus).
				WithField("kind", e.Kind).
				WithField("developer_message", e.DeveloperMessage)
			if e.HTTPStatus == http.StatusForbidden || e.HTTPStatus == http.StatusUnauthorized {
				l.Warnf("api: api token access error. %s", e.Error())
			} else {
				l.Debugf("api: api generic error. %s", e.Error())
			}
			c.JSON(e.HTTPStatus, e)
		case *errors.InputDataError:
			logger.
				WithField("http_status", e.HTTPStatus).
				WithField("kind", e.Kind).
				WithField("developer_message", e.DeveloperMessage).
				Debugf("api: api generic error. %s", e.Error())
			c.JSON(e.HTTPStatus, e)
		case *echo.HTTPError:
			logger.WithError(e).WithField("errorCode", e.Code).Errorf("api: echo http error. %s", e.Error())
			c.JSON(e.Code, e)
		default:
			apiErr := errors.NewInternalServerError()
			if config.Debug {
				apiErr.Err = err.Error()
			}
			logger.WithError(err).Error("api: internal server error.")
			c.JSON(apiErr.HTTPStatus, apiErr)
		}
	}
	// CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.Cors.AllowOrigins,
		AllowHeaders:     config.Cors.AllowHeaders,
		ExposeHeaders:    config.Cors.ExposeHeaders,
		AllowCredentials: config.Cors.AllowCredentials,
	}))

	apirouter := router.Group("/v1")

	// Create custom context for the user and also read the token
	apirouter.Use(middlewares.RequestContextMiddleware(logger, helpers.Token.DecodeToken))

	// Serve resources
	handlers.ServeProductResource(apirouter)

	return router
}
