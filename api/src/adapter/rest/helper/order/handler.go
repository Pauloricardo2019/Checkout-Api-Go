package helper

import (
	"errors"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	model "ravxcheckout/src/internal/model/db"

	"github.com/google/uuid"
)

type GetExistOrderFn func(
	orderID string,
	fullOrder bool,
	payed bool,
	getOrderById order.GetByIDFn,
) (*model.Order, int, error)

func GetExistOrder(
	orderID string,
	fullOrder bool,
	payed bool,
	getOrderById order.GetByIDFn,
) (*model.Order, int, error) {

	if orderID == "" {
		return nil, 400, errors.New("\"order_id\" is not allowed to be empty")
	}

	_, err := uuid.Parse(orderID)
	if err != nil {
		return nil, 400, errors.New("\"order_id\" error: " + err.Error())
	}

	order, err := getOrderById(orderID, !fullOrder, payed)
	if err != nil {
		statusCode := 404
		if err.Error() != "record not found" {
			statusCode = 500
		}
		return nil, statusCode, err
	}

	return order, 200, nil
}
