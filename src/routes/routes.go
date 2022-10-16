package routes

import (
	"net/http"

	ctl "github.com/silastgoes/mock-store/src/controllers"
)

var (
	templatePath = "templates/*.html"
)

type router struct {
	pcs ctl.ProductControlService
}

//go:generate mockgen --source=routes.go --package=mocks --destination=./mocks/routes.go  RouterService
type RouterService interface {
	LoadRoutes()
}

func NewRouterService(controller ctl.ProductControlService) *router {
	return &router{
		pcs: controller,
	}
}

func (r *router) LoadRoutes() {
	http.HandleFunc("/", r.pcs.Index)
	http.HandleFunc("/new", r.pcs.New)
	http.HandleFunc("/insert", r.pcs.Insert)
	http.HandleFunc("/delete", r.pcs.Delete)
	http.HandleFunc("/edit", r.pcs.Edit)
	http.HandleFunc("/update", r.pcs.Update)
}
