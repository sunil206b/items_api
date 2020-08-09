package repo

import (
	"database/sql"
	"github.com/sunil206b/items_api/model"
	"github.com/sunil206b/store_utils_go/errors"
	"net/http"
)

type IItemsRepo interface {
	CreateItem(item *model.Item) *errors.RestErr
	GetItem(itemId string) (*model.Item, *errors.RestErr)
}

type itemRepo struct {
	conn *sql.DB
}

func NewItemRepo(db *sql.DB) IItemsRepo {
	return &itemRepo{
		conn: db,
	}
}

func (rep *itemRepo) CreateItem(item *model.Item) *errors.RestErr {
	return &errors.RestErr{
		Message: "Implement me!",
		StatusCode: http.StatusNotImplemented,
		Error: "not_implemented",
	}
}

func (rep *itemRepo) GetItem(itemId string) (*model.Item, *errors.RestErr) {
	return nil, &errors.RestErr{
		Message: "Implement me!",
		StatusCode: http.StatusNotImplemented,
		Error: "not_implemented",
	}
}