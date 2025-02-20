package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"backend/routes"
	"backend/utils"
)

func main() {
	utils.InitDB()

	router := mux.NewRouter()
	routes.InitRoutes(router)

	log.Println(http.ListenAndServe(":8000", router))
}
