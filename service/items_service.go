package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sunil206b/items_api/model"
	"github.com/sunil206b/items_api/repo"
	"github.com/sunil206b/store_utils_go/errors"
	"github.com/sunil206b/store_utils_go/logger"
)

const (
	indexItems = "items"
)

type IItemService interface {
	CreateItem(item *model.Item) *errors.RestErr
	GetItem(itemId string) (*model.Item, *errors.RestErr)
	Search(query model.EsQuery) ([]*model.Item, *errors.RestErr)
	Delete(itemId string) (string, *errors.RestErr)
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
	result, err := srv.getIndex(indexItems, itemId)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error whe trying to get id %s", itemId))
	}
	if !result.Found {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no item found with id %s", itemId))
	}
	var item model.Item
	rawData, err := result.Source.MarshalJSON();
	if err != nil {
		return nil, errors.NewInternalServerError("error when trying to parse data from database")
	}
	err = json.Unmarshal(rawData, &item)
	if err != nil {
		return nil, errors.NewInternalServerError("error when trying to parse data from database")
	}
	item.ItemId = itemId
	return &item, nil
}

func (srv *itemService) Search(query model.EsQuery) ([]*model.Item, *errors.RestErr) {
	result, err := srv.searchItems(indexItems, query.Build())
	if err != nil {
		return nil, errors.NewInternalServerError("error when trying to search documents")
	}
	items := make([]*model.Item, 0, result.TotalHits())
	for _, hit := range result.Hits.Hits {
		var item model.Item
		bytes, _ := hit.Source.MarshalJSON()
		if err = json.Unmarshal(bytes, &item); err != nil {
			return nil, errors.NewInternalServerError("error when trying to parse response")
		}
		item.ItemId = hit.Id
		items = append(items, &item)
	}
	if len(items) == 0 {
		return nil, errors.NewNotFoundError("no items found with given criteria")
	}
	return items, nil
}

func (srv *itemService) Delete(itemId string) (string, *errors.RestErr) {
	res, err := srv.delete(indexItems, itemId)
	if err != nil {
		return "not deleted", errors.NewNotFoundError(fmt.Sprintf("error when delete the item id %s", itemId))
	}
	return res.Result, nil
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

func (srv *itemService) getIndex(index string, itemId string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := srv.client.Get().Index(index).Id(itemId).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to get item by id %s", itemId), err)
		return nil, err
	}
	return result, err
}

func (srv *itemService) searchItems(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	result, err := srv.client.Search(index).Query(query).RestTotalHitsAsInt(true).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (srv *itemService) delete(index string, itemId string) (*elastic.DeleteResponse, error) {
	ctx := context.Background()
	res, err := srv.client.Delete().Index(index).Type("_doc").Id(itemId).Pretty(true).Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}