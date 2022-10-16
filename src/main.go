package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/silastgoes/mock-store/src/controllers"
	"github.com/silastgoes/mock-store/src/dbconnection"
	"github.com/silastgoes/mock-store/src/model/product"

	rts "github.com/silastgoes/mock-store/src/routes"
)

var (
	templatePath = "templates/*.html"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func main() {
	db, _ := dbconnection.NewDatabadeConnection().GetDb()
	defer db.Close()
	LoadControlles(db)
	log.Fatal(http.ListenAndServe(":4444", nil))
}

func LoadControlles(db *sql.DB) {
	srv := product.NewProductModelService(db)
	pc := controllers.NewProductControl(templatePath, srv)
	rts.NewRouterService(pc).LoadRoutes()
}
