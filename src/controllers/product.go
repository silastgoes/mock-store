package controllers

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/silastgoes/mock-store/src/model/product"
)

type productControl struct {
	productService product.ProductModelService
	Template       *template.Template
}

//go:generate mockgen --source=product.go --package=mocks --destination=./mocks/product.go  ProductControlService
type ProductControlService interface {
	Index(w http.ResponseWriter, r *http.Request)
	New(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
}

func NewProductControl(path string, svr product.ProductModelService) *productControl {
	temp := template.Must(template.ParseGlob(path))

	return &productControl{
		productService: svr,
		Template:       temp,
	}
}

func (pc *productControl) Index(w http.ResponseWriter, r *http.Request) {
	products, err := pc.productService.GetProducts()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Erro em recuperação de produtos:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	pc.Template.ExecuteTemplate(w, "Index", products)
}

func (pc *productControl) New(w http.ResponseWriter, r *http.Request) {
	pc.Template.ExecuteTemplate(w, "NewProduct", nil)
}

func (pc *productControl) Insert(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMovedPermanently
	if r.Method == "POST" {
		name := r.FormValue("name")
		description := r.FormValue("description")
		value := r.FormValue("value")
		quantity := r.FormValue("quantity")

		convertedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			status = http.StatusBadRequest
			log.Println("Erro na converção de preço:", err)
		}

		convertedQuantity, err := strconv.Atoi(quantity)
		if err != nil {
			status = http.StatusBadRequest
			log.Println("Erro na converção de quantidade:", err)
		}

		if status == http.StatusMovedPermanently {
			err = pc.productService.Create(name, description, convertedValue, convertedQuantity)
			if err != nil {
				log.Println("Erro na criação de produto:", err)
				status = http.StatusInternalServerError
			}
		}
	}

	http.Redirect(w, r, "/", status)
}

func (pc *productControl) Delete(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMovedPermanently
	id := r.URL.Query().Get("id")

	err := pc.productService.Delete(id)
	if err != nil {
		log.Println("Erro ao deletar um produto:", err)
		status = http.StatusInternalServerError
	}

	http.Redirect(w, r, "/", status)
}

func (pc *productControl) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := http.StatusOK

	p, err := pc.productService.Get(id)
	if err != nil {
		log.Println("Erro na busca de produtos:", err)
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	pc.Template.ExecuteTemplate(w, "Edit", p)
}

func (pc *productControl) Update(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMovedPermanently
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		description := r.FormValue("description")
		value := r.FormValue("value")
		quantity := r.FormValue("quantity")

		convertedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println("Erro na converção de preço:", err)
			status = http.StatusBadRequest
		}

		convertedQuantity, err := strconv.Atoi(quantity)
		if err != nil {
			log.Println("Erro na converção de quantidade:", err)
			status = http.StatusBadRequest
		}

		convertedId, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Erro na converção de id:", err)
			status = http.StatusNotFound
		}

		if status == http.StatusMovedPermanently {
			err = pc.productService.Update(convertedId, name, description, convertedValue, convertedQuantity)
			if err != nil {
				log.Println("Erro no update de produto:", err)
				status = http.StatusInternalServerError
			}
		}
	}

	http.Redirect(w, r, "/", status)
}
