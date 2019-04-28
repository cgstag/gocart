package datastores

import (
	"fmt"

	_ "github.com/Kount/pq-timeouts"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Database DatastoreInterface

type DatastoreInterface interface {
	SetLogger(logger logger)
	Close() error
	BeginTransaction() (tx *gorm.DB)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
	ProductsGetter
}

type logger interface {
	Print(v ...interface{})
}

type datastore struct {
	db *gorm.DB
}

func SetupDatastore(dialect, dsn string, maxIdleConnections, maxOpenConnections int, debug bool) error {
	db, err := gorm.Open(string(dialect), dsn)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("error when initializing the datastore. %s", err)
	}
	// Debug mode on
	if debug {
		db.LogMode(true)
	}
	// Setup max idle and open connections
	db.DB().SetMaxIdleConns(maxIdleConnections)
	db.DB().SetMaxOpenConns(maxOpenConnections)

	// Execute ping against database to check connection
	err = db.DB().Ping()
	if err != nil {
		return fmt.Errorf("dsite datastore: error when pinging the database. %s", err)
	}
	Database = &datastore{db}
	return nil
}

func (s *datastore) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("dsite datastore: error when trying to close the datastore. %s", err)
	}
	return nil
}

func (s *datastore) SetLogger(logger logger) {
	s.db.SetLogger(logger)
}

func (s *datastore) BeginTransaction() (tx *gorm.DB) {
	return s.db.Begin()
}

func (s *datastore) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (s *datastore) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (s *datastore) Products() ProductDatastoreInterface {
	return newProductDatastore(s.db)
}
