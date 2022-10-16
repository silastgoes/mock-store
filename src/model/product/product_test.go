package product

import (
	"fmt"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/silastgoes/mock-store/src/util"
	"github.com/stretchr/testify/assert"
)

// RandonProduct generate a random product
func RandonProduct() Product {
	return Product{
		Id:          util.RandomInt(5, 100),
		Name:        util.RandomString(6),
		Description: util.RandomString(20),
		Quantity:    util.RandomInt(1, 2000),
		Value:       util.RandomFloat(),
	}
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(err)

	rows := &sqlmock.Rows{}
	coluns := []string{"id", "name", "description", "value", "quantity"}
	result := RandonProduct()
	ps := NewProductModelService(db)

	t.Run("Testing success result", func(t *testing.T) {

		rows = sqlmock.NewRows(coluns).
			AddRow(
				result.Id,
				result.Name,
				result.Description,
				result.Value,
				result.Quantity,
			)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM product WHERE id = $1`)).
			WithArgs(fmt.Sprint(result.Id)).
			WillReturnRows(rows)

		res, err := ps.Get(fmt.Sprint(result.Id))

		assert.Nil(err)
		assert.Equal(res.Id, result.Id)
		assert.Equal(res.Name, result.Name)
		assert.Equal(res.Description, result.Description)
		assert.Equal(res.Quantity, result.Quantity)
		assert.Equal(res.Value, result.Value)
	})

	t.Run("Testing error Query", func(t *testing.T) {
		rows = sqlmock.NewRows(coluns).
			AddRow(
				"1",
				result.Name,
				result.Description,
				result.Value,
				result.Quantity,
			)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM product WHERE id = $1`)).
			WithArgs(2).
			WillReturnRows(rows)

		_, err := ps.Get(fmt.Sprint(result.Id))

		assert.Error(err)
	})

	t.Run("Testing error scan", func(t *testing.T) {
		rows = sqlmock.NewRows(coluns).
			AddRow(
				1.5,
				result.Name,
				result.Description,
				result.Value,
				result.Quantity,
			)

		mock.ExpectQuery(`SELECT * FROM product WHERE id = $1`).
			WithArgs(1.5).
			WillReturnRows(rows)

		_, err := ps.Get(fmt.Sprint(result.Id))

		assert.Error(err)
	})
}

func TestGetProducts(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(err)

	rows := &sqlmock.Rows{}
	coluns := []string{"id", "name", "description", "value", "quantity"}
	result := RandonProduct()
	ps := NewProductModelService(db)

	t.Run("Testing success result", func(t *testing.T) {

		rows = sqlmock.NewRows(coluns).
			AddRow(
				result.Id,
				result.Name,
				result.Description,
				result.Value,
				result.Quantity,
			)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM product ORDER BY id ASC`)).
			WithArgs().
			WillReturnRows(rows)

		res, err := ps.GetProducts()

		assert.Nil(err)
		assert.Equal(res[0].Id, result.Id)
		assert.Equal(res[0].Name, result.Name)
		assert.Equal(res[0].Description, result.Description)
		assert.Equal(res[0].Quantity, result.Quantity)
		assert.Equal(res[0].Value, result.Value)
	})

	t.Run("Testing error Query", func(t *testing.T) {
		rows = sqlmock.NewRows(coluns).
			AddRow(
				"1",
				result.Name,
				result.Description,
				result.Value,
				result.Quantity,
			)

		mock.ExpectQuery(`SELECT * FROM product ORDER BY id ASC`).
			WithArgs().
			WillReturnRows(rows)

		_, err := ps.GetProducts()

		assert.Error(err)
	})
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(err)

	result := RandonProduct()
	ps := NewProductModelService(db)

	t.Run("Testing success result", func(t *testing.T) {

		prepare := regexp.QuoteMeta("INSERT INTO product(name, description, value, quantity) VALUES($1, $2, $3, $4)")
		mock.ExpectPrepare(prepare).
			ExpectExec().
			WithArgs(result.Name, result.Description, result.Value, result.Quantity)

		err := ps.Create(result.Name, result.Description, result.Value, result.Quantity)

		assert.Nil(err)
	})

	t.Run("Testing Error", func(t *testing.T) {

		prepare := "INSERT INTO product(name, description, value, quantity) VALUES($1, $2, $3, $4)"
		mock.ExpectPrepare(prepare).
			ExpectExec().
			WithArgs(result.Name, result.Description, result.Value, result.Quantity)

		err := ps.Create(result.Name, result.Description, result.Value, result.Quantity)

		assert.Error(err)

	})
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(err)

	result := RandonProduct()
	ps := NewProductModelService(db)

	t.Run("Testing success result", func(t *testing.T) {

		prepare := regexp.QuoteMeta("UPDATE product SET name=$1 , description=$2, value=$3, quantity=$4 WHERE id=$5")
		mock.ExpectPrepare(prepare).
			ExpectExec().
			WithArgs(result.Name, result.Description, result.Value, result.Quantity, result.Id)

		err := ps.Update(result.Id, result.Name, result.Description, result.Value, result.Quantity)

		assert.Nil(err)
	})

	t.Run("Testing Error", func(t *testing.T) {

		prepare := "UPDATE product SET name=$1 , description=$2, value=$3, quantity=$4 WHERE id=$5"
		mock.ExpectPrepare(prepare).
			ExpectExec().
			WithArgs(result.Name, result.Description, result.Value, result.Quantity, result.Id)

		err := ps.Update(result.Id, result.Name, result.Description, result.Value, result.Quantity)

		assert.Error(err)

	})
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(err)

	result := RandonProduct()
	ps := NewProductModelService(db)

	t.Run("Testing success result", func(t *testing.T) {

		prepare := regexp.QuoteMeta("DELETE FROM product WHERE id=$1")
		mock.ExpectPrepare(prepare).
			ExpectExec().
			WithArgs(fmt.Sprint(result.Id))

		err := ps.Delete(fmt.Sprint(result.Id))

		assert.Nil(err)
	})

	t.Run("Testing Error", func(t *testing.T) {

		prepare := "DELETE FROM product WHERE id=$1"
		mock.ExpectPrepare(prepare).
			ExpectExec().
			WithArgs(fmt.Sprint(result.Id))

		err := ps.Delete(fmt.Sprint(result.Id))

		assert.Error(err)

	})
}
