package postgresql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	address  = "localhost:5432"
	username = "golang"
	password = ""
	database = "godb"
)

type User struct {
	ID       int
	Username string
	Name     string
}

func TestNew(t *testing.T) {
	db, err := New(address, username, password, database)
	if err != nil {
		fmt.Println("PostgreSQL failed")
	}

	assert.NotNil(t, db)
	assert.Nil(t, err)
}

func TestCreateSchema(t *testing.T) {
	db, _ := New(address, username, password, database)
	modelsList := []interface{}{&User{}}

	// Should be able to create schema
	err := CreateSchema(db, modelsList, true)
	assert.Nil(t, err)

	// Should be able to query schema
	var info []struct {
		ColumnName string
		DataType   string
	}

	_, err = db.Query(&info, `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = 'users'`)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(info)

	assert.Nil(t, err)
	assert.Equal(t, info[0].ColumnName, "id")
	assert.Equal(t, info[0].DataType, "bigint")
	assert.Equal(t, info[1].ColumnName, "username")
	assert.Equal(t, info[1].DataType, "text")
	assert.Equal(t, info[2].ColumnName, "name")
	assert.Equal(t, info[2].DataType, "text")
}
