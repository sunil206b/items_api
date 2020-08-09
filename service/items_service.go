package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sunil206b/items_api/logger"
	"github.com/sunil206b/items_api/model"
	"github.com/sunil206b/items_api/repo"
	"github.com/sunil206b/store_utils_go/errors"
	"src/github.com/olivere/elastic"
)

const (
	indexItems = "items"
)

type IItemService interface {
	CreateItem(item *model.Item) *errors.RestErr
	GetItem(itemId string) (*model.Item, *errors.RestErr)
}

type itemService struct {
	repo repo.IItemsRepo
	client *elastic.Client
}

func NewItemService(db *sql.DB, esClient *elastic.Client) IItemService {
	return &itemService{
		repo: repo.NewItemRepo(db),
		client: esClient,
	}
}

func (srv *itemService) CreateItem(item *model.Item) *errors.RestErr {
	result, err := srv.index(indexItems, item)
	if err != nil {
		return errors.NewInternalServerError("error when trying to save item")
	}
	item.ItemId = result.Id
	return nil
}

func (srv *itemService) GetItem(itemId string) (*model.Item, *errors.RestErr) {
	return srv.repo.GetItem(itemId)
}

func (srv *itemService) index(index string, item interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err :=  srv.client.Index().Index(index).BodyJson(item).Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to index document in index %s", index), err)
		return nil, err
	}
	return result, nil
}