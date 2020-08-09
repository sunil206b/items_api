package app

import (
	"database/sql"
	"github.com/sunil206b/items_api/controller"
	"github.com/sunil206b/items_api/driver"
	"net/http"
)

func mapUrls(db *sql.DB) {
	itemH := controller.NewItemController(db, driver.GetEsClient())
	router.HandleFunc("/ping", itemH.Ping)
	router.HandleFunc("/items", itemH.CreateItem).Methods(http.MethodPost)
	router.HandleFunc("/items/{item_id}", itemH.GetById)
}
