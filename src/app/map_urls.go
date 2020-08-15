package app

import (
	"database/sql"
	"github.com/sunil206b/items_api/src/controller"
	"github.com/sunil206b/items_api/src/driver"
	"net/http"
)

func mapUrls(db *sql.DB) {
	itemH := controller.NewItemController(db, driver.GetEsClient())
	router.HandleFunc("/ping", itemH.Ping)
	router.HandleFunc("/items", itemH.CreateItem).Methods(http.MethodPost)
	router.HandleFunc("/items/{item_id}", itemH.GetById).Methods(http.MethodGet)
	router.HandleFunc("/items/search", itemH.Search).Methods(http.MethodPost)
	router.HandleFunc("/items/{item_id}", itemH.Delete).Methods(http.MethodDelete)
}
