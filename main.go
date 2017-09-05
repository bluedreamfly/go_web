package main

import (

	"log"
	"net/http"

	"./controllers"
	"./utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

type httpServer struct {

}

func (server httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request){


	route := utils.Route{}

	db, err := gorm.Open("mysql", "root:@/ebus?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {

		println(err)
	}

	route.Get("/products", func(w http.ResponseWriter, r *utils.Req) {
		controllers.GetList(w, r, db)
	})

	route.Post("/product", func(w http.ResponseWriter, r *utils.Req) {
		controllers.AddProduct(w, r, db)
	})

	route.PUT("/product/:id", func(w http.ResponseWriter, r *utils.Req) {
		controllers.UpdateProduct(w, r, db)
	})

	route.DELETE("/product/:id", func(w http.ResponseWriter, r *utils.Req) {
		controllers.DelProduct(w, r, db)
	})



	route.Start(w, r)
}

func main() {
	var server httpServer
	http.Handle("/", server)

	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		log.Fatal("listenandserve", err)
	}
}
