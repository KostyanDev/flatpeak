package http

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	_ "app/docs"
)

func RegisterRoutes(router *mux.Router, handler *Handler) {
	router.HandleFunc("/slots", handler.GetOptimalSlots).Methods(http.MethodGet)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
