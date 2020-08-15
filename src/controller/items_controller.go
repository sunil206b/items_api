package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/sunil206b/items_api/src/model"
	"github.com/sunil206b/items_api/src/service"
	"github.com/sunil206b/oauth_go/oauth"
	"github.com/sunil206b/store_utils_go/errors"
	"github.com/sunil206b/store_utils_go/http_utils"
	"github.com/sunil206b/store_utils_go/logger"
	"net/http"
	"strings"
)

type itemController struct {
	srv service.IItemService
}

func NewItemController(db *sql.DB, esClient *elastic.Client) *itemController {
	return &itemController{
		srv: service.NewItemService(db, esClient),
	}
}

func (c *itemController) Ping(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("I am running"))
}

func (c *itemController) CreateItem(res http.ResponseWriter, req *http.Request) {
	if err := oauth.AuthenticateRequest(req); err != nil {
		http_utils.ResponseError(res, err)
		return
	}

	var item model.Item
	err := json.NewDecoder(req.Body).Decode(&item)
	if err != nil {
		logger.Error("error when trying to marshal item", err)
		errMsg := errors.NewBadRequest("Not a valid Item json")
		http_utils.ResponseError(res, errMsg)
		return
	}
	item.Seller = oauth.GetCallerId(req)
	if errMsg := c.srv.CreateItem(&item); errMsg != nil {
		http_utils.ResponseError(res, errMsg)
		return
	}
	http_utils.ResponseJson(res, http.StatusCreated, item)
}

func (c *itemController) GetById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	itemId := strings.TrimSpace(vars["item_id"])
	item, err := c.srv.GetItem(itemId)
	if err != nil {
		http_utils.ResponseError(res, err)
		return
	}
	http_utils.ResponseJson(res, http.StatusOK, item)
}

func (c *itemController) Search(res http.ResponseWriter, req *http.Request) {
	var query model.EsQuery
	err := json.NewDecoder(req.Body).Decode(&query)
	if err != nil {
		errMsg := errors.NewBadRequest("invalid json body")
		http_utils.ResponseError(res, errMsg)
		return
	}
	items, errMsg := c.srv.Search(query)
	if errMsg != nil {
		http_utils.ResponseError(res, errMsg)
		return
	}
	http_utils.ResponseJson(res, http.StatusOK, items)
}

func (c *itemController) Delete(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	itemId := strings.TrimSpace(vars["item_id"])
	if itemId == "" {
		errMsg := errors.NewBadRequest("not a valid item id")
		http_utils.ResponseError(res, errMsg)
		return
	}
	msg, errMsg := c.srv.Delete(itemId)
	if errMsg != nil {
		http_utils.ResponseError(res, errMsg)
		return
	}
	http_utils.ResponseJson(res, http.StatusOK, map[string]string{"item status": msg})
}