package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	mux "github.com/gorilla/mux"
	ents "github.com/sabrs0/L0_WB/entities"
)

type OrderHandler struct {
	ordersCache *ents.OrdersCache
}

func (h OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		fmt.Println("id is ", id)
	}
	ords := (*h.ordersCache)[id]
	data, err := json.Marshal(ords)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Println("http data is", data)
		w.Write(data)
	}

}

func StartServer(addr string, cache *ents.OrdersCache) {

	router := mux.NewRouter()

	router.Handle("/orders/", OrderHandler{ordersCache: cache})

	server := http.Server{
		Addr: addr,

		Handler: router,
	}
	fmt.Println("Start listening at ", server.Addr)
	server.ListenAndServe()
}
