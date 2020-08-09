package service

import (
	"context"
	"database/sql"
	"github.com/sunil206b/items_api/model"
	"github.com/sunil206b/items_api/repo"
	"github.com/sunil206b/store_utils_go/errors"
	"src/github.com/olivere/elastic"
)

type IItemService interface {
	CreateItem(item *model.Item) *errors.RestErr
	GetItem(itemId string) (*model.Item, *errors.RestErr)
	Index(interface{}) (*elastic.IndexResponse, error)
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
	return srv.repo.CreateItem(item)
}

func (srv *itemService) GetItem(itemId string) (*model.Item, *errors.RestErr) {
	return srv.repo.GetItem(itemId)
}

func (srv *itemService) Index(interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	return srv.client.Index().Do(ctx)
}