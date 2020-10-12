package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	helper "modulorgo/helpers"
	route "modulorgo/routes"
)

func main() {
	if os.Getenv("JWT_KEY") == "" {
		log.Print("JWT_KEY envar is null, please set JWT_KEY in enviroment ")
		return
	}

	db := helper.InitDB()
	defer db.Close()
	router := mux.NewRouter()
	route.Register(db, router)
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			log.Print("Error when registered Path")
			return nil
		}
		method, err := route.GetMethods()
		if err != nil {
			log.Print("Error when registered Method")
			return nil
		}
		log.Print("Endpoint [" + string(method[0]) + ":" + path + "] Registered")
		return nil
	})
	log.Print("Server running on port ", os.Getenv("BACKEND_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("BACKEND_PORT"), handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Cache-Control"}),
		handlers.ExposedHeaders([]string{"Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router),
	))
}
