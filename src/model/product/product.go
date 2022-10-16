package product

import (
	"database/sql"
)

type Product struct {
	Id          int
	Name        string
	Description string
	Value       float64
	Quantity    int
}

type productModel struct {
	DB *sql.DB
}

//go:generate mockgen --source=product.go --package=mocks --destination=./mocks/product.go  ProductService
type ProductModelService interface {
	Create(name, description string, value float64, quantity int) error
	Get(param string) (Product, error)
	GetProducts() ([]Product, error)
	Update(id int, name, description string, value float64, quantity int) error
	Delete(id string) error
}

func NewProductModelService(db *sql.DB) *productModel {
	return &productModel{
		DB: db,
	}
}

func (prod *productModel) GetProducts() (products []Product, err error) {
	p := Product{}

	rows, err := prod.DB.Query("SELECT * FROM product ORDER BY id ASC")
	if err != nil {
		return
	}

	for rows.Next() {
		var id, quantity int
		var name, description string
		var value float64

		err = rows.Scan(&id, &name, &description, &value, &quantity)
		if err != nil {
			return
		}

		p.Id = id
		p.Name = name
		p.Description = description
		p.Value = value
		p.Quantity = quantity

		products = append(products, p)
	}

	return
}

func (prod *productModel) Update(
	id int,
	name, description string,
	value float64,
	quantity int,

) error {
	rows, err := prod.DB.Prepare("UPDATE product SET name=$1 , description=$2, value=$3, quantity=$4 WHERE id=$5")
	if err != nil {
		return err
	}

	rows.Exec(name, description, value, quantity, id)
	return nil
}

func (prod *productModel) Create(name, description string, value float64, quantity int) error {

	rows, err := prod.DB.Prepare("INSERT INTO product(name, description, value, quantity) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	rows.Exec(name, description, value, quantity)
	return nil
}

func (prod *productModel) Delete(id string) error {
	rows, err := prod.DB.Prepare("DELETE FROM product WHERE id=$1")
	if err != nil {
		return err
	}

	rows.Exec(id)
	return nil
}

func (prod *productModel) Get(param string) (Product, error) {
	p := Product{}

	rows, err := prod.DB.Query("SELECT * FROM product WHERE id = $1", param)
	if err != nil {
		return p, err
	}

	for rows.Next() {
		var id, quantity int
		var name, description string
		var value float64

		err = rows.Scan(&id, &name, &description, &value, &quantity)
		if err != nil {
			return p, err
		}

		p.Id = id
		p.Name = name
		p.Description = description
		p.Value = value
		p.Quantity = quantity
	}

	return p, nil
}
