package route

import (
	"github.com/gorilla/mux"
	"mocking_server/controllers"
)

type Router struct {
}

func (ro *Router) Make() *mux.Router {

	root := mux.NewRouter()

	v3 := root.PathPrefix("/v3").Subrouter()

	narindo := new(controllers.NarindoController)

	v3.HandleFunc("/h2h", narindo.TopUp).Methods("POST")
	v3.HandleFunc("/advice", narindo.CheckStatus).Methods("POST")
	v3.HandleFunc("/change-status", narindo.ChangeStatus).Methods("POST")
	v3.HandleFunc("/change-credential", narindo.ChangeCredential).Methods("POST")

	return v3
}
