package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/sunil206b/items_api/logger"
	"github.com/sunil206b/items_api/model"
	"github.com/sunil206b/items_api/service"
	"github.com/sunil206b/oauth_go/oauth"
	"github.com/sunil206b/store_utils_go/errors"
	"github.com/sunil206b/store_utils_go/http_utils"
	"net/http"
	"src/github.com/olivere/elastic"
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

}