package handlers

import (
	"fmt"
	"strconv"

	"net/http"

	"github.com/labstack/echo"
	cache "github.com/patrickmn/go-cache"
	"api/errors"
	"api/middlewares"
	"api/models"
	"api/stores/datastores"
)

type productResource struct {
	cache *cache.Cache
}

func ServeProductResource(router *echo.Group) {
	r := &productResource{Cache}

	rg := router.Group("/product")

	rg.GET("/:id", r.getProductById)
}

/**
GetBaseProductByIdProduct : return a Product
@path	/product/:id
@return JSON
*/
func (r *productResource) getProductById(c echo.Context) error {
	// Get Logger
	rc := middlewares.GetRequestContext(c)
	logger := rc.Log
	// Get Parameter
	id := c.Param("id")
	// Check if the given ID is numerical
	idProduct, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewInputDataError(fmt.Errorf("invalid product ID %s, must be an integer", c.Param("id")))
	}

	// First check if a cache entry exists, return it if found
	if cacheValue, found := r.cache.Get("get-product-by-id_" + id); found {
		logger.WithField("id", id).Info("getProductById() : Cache Found")

		return c.JSON(200, cacheValue.(*models.Product))
	} else {
		logger.WithField("id", id).Debug("getSkuByProductId() : Retrieving Product")
		result, err := datastores.Database.Products().GetOne(idProduct)

		if err != nil {
			// Cache the Request with no sensible result
			logger.WithField("id", id).Info("No Product found for this ID")
			return c.JSON(http.StatusAccepted, err)
		}

		// Set a Cache and return de result
		r.cache.Set("get-product-by-id_"+id, result, cache.DefaultExpiration)
		logger.WithField("id", id).WithField("baseProduct", result.Name).Info("GetProductById OK (http/200)")
		return c.JSON(200, result)
	}
}
