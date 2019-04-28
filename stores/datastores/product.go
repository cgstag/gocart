package datastores

import (
	"github.com/jinzhu/gorm"
	"./errors"
	"./models"
)

type productDatastore struct {
	db *gorm.DB
}

type ProductsGetter interface {
	Products() ProductDatastoreInterface
}

type ProductDatastoreInterface interface {
	GetOne(id int) (*models.Product, error)
}

func newProductDatastore(db *gorm.DB) ProductDatastoreInterface {
	return &productDatastore{db}
}

/**
// TODO Go Docs
*/
func (s *productDatastore) GetOne(id int) (*models.Product, error) {
	p := &models.Product{}
	if r := s.db.
		First(p, id); r.Error != nil {
		if r.RecordNotFound() {
			return nil, errors.NewRecordNotFound()
		}
		return nil, r.Error
	}

	return p, nil
}
