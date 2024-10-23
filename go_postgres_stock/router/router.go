package router

import (
	"github.com/gorilla/mux"
	"github.com/kaung-minkhant/go_projs/go_postgres_stock/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/stocks/{id}", controllers.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stocks", controllers.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stocks", controllers.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stocks/{id}", controllers.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/stocks/{id}", controllers.DeleteStock).Methods("DELETE", "OPTIONS")

	return router
}
