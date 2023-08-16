package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	mux "github.com/gorilla/mux"
	ents "github.com/sabrs0/L0_WB/entities"
)

type OrderHandler struct {
	ordersCache *ents.OrdersCache
}

func (h OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../../templates/orders.html"))

	id := r.URL.Query().Get("id")
	if id != "" {
		fmt.Println("id is ", id)
	}
	ords := []ents.Orders{(*h.ordersCache)[id]}
	_, err := json.Marshal(ords)
	if err != nil {
		http.Error(w, err.Error(), 404)
	} else {

		tmpl.Execute(w, struct {
			Orders []ents.Orders
		}{ords})
	}

}

func StartServer(addr string, cache *ents.OrdersCache) {

	router := mux.NewRouter()

	router.Handle("/orders", OrderHandler{ordersCache: cache})

	server := http.Server{
		Addr: addr,

		Handler: router,
	}
	fmt.Println("Start listening at ", server.Addr)
	server.ListenAndServe()
}
