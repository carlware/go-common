package postgresql

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

const (
	tableExistsErrorCode = "ERROR #42P07"
)

// New returns a new Postgres Database instance
func New(address string, user string, password string, database string) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     address,
		User:     user,
		Password: password,
		Database: database,
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		fmt.Println(err)
		db.Close()
		return nil, err
	}

	return db, nil
}

// CreateSchema creates the database tables if dropExisting is set to true it will drop the current schema
func CreateSchema(db *pg.DB, modelsList []interface{}, dropExisting bool) error {
	for _, model := range modelsList {
		if dropExisting {
			err := DropTable(db, model)
			if err != nil {
				fmt.Println(err)
			}
		}
		err := CreateTable(db, model)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTable creates a new table in the database
func CreateTable(db *pg.DB, model interface{}) error {
	err := db.CreateTable(model, &orm.CreateTableOptions{
		Temp: false,
	})
	if err != nil {
		errorCode := err.Error()[0:12]
		if errorCode != tableExistsErrorCode {
			return err
		}
	}

	return nil
}

// DropTable deletes the existing tables
func DropTable(db *pg.DB, model interface{}) error {
	err := db.DropTable(model, &orm.DropTableOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Ping test connection
func Ping(db *pg.DB) error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		fmt.Println(err)
		db.Close()
		return err
	}
	return nil
}
