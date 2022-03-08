package notificationcontroller

import (
	"ravxcheckout/src/adapter/database/sql/repository/order"
	modelDB "ravxcheckout/src/internal/model/db"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetByID_OK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		assert.False(t, fullOrder)
		assert.True(t, payed)

		result := &modelDB.Order{
			ID:         "267511c9-324a-41dd-8711-c7d72e179522",
			Status:     "WAITING_PAYMENT",
			Amount:     1000,
			Currency:   "BRL",
			CustomerID: "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
		}
		return result, 200, nil
	}

	getQueryFromParamsMock := func(context *gin.Context, params string) string {
		paramsMock := map[string]string{
			"order_id":    "267511c9-324a-41dd-8711-c7d72e179522",
			"status":      "APPROVED",
			"customer_id": "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			"payment_id":  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}

		param, ok := paramsMock[params]
		assert.True(t, ok, "Provided 'params' is invalid")
		assert.NotEmpty(t, param)

		return param
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		status, ok := changes["status"]
		assert.True(t, ok, "'status' need to be provided")
		assert.Equal(t, "APPROVED", status)
		payment_id, ok := changes["payment_id"]
		assert.True(t, ok, "'payment_id' need to be provided")
		assert.Equal(t, "b10763c4-9f63-42d8-83ca-40824734ec6d", payment_id)
		return nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 200, statusCode)
		assert.Nil(t, obj)
	}

	returnErrorMock := func(context *gin.Context, statusCode int, err error) {
		assert.Fail(t, "An error happened, when is not expected error: "+err.Error())
	}

	NotificationHandler(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		getQueryFromParamsMock,
		returnJSONMock,
		returnErrorMock,
	)(nil)
}

func TestGetByID_AlreadyPayed_OK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		assert.False(t, fullOrder)
		assert.True(t, payed)

		pay_id := "b10763c4-9f63-42d8-83ca-40824734ec6d"
		result := &modelDB.Order{
			ID:         "267511c9-324a-41dd-8711-c7d72e179522",
			Status:     "APPROVED",
			Amount:     1000,
			Currency:   "BRL",
			CustomerID: "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			PaymentID:  &pay_id,
		}
		return result, 200, nil
	}

	getQueryFromParamsMock := func(context *gin.Context, params string) string {
		paramsMock := map[string]string{
			"order_id":    "267511c9-324a-41dd-8711-c7d72e179522",
			"status":      "APPROVED",
			"customer_id": "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			"payment_id":  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}

		param, ok := paramsMock[params]
		assert.True(t, ok, "Provided 'params' is invalid")
		assert.NotEmpty(t, param)

		return param
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		assert.Fail(t, "An error happened, 'updateStatus' is not to be called")
		return nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 202, statusCode)
		assert.Nil(t, obj)
	}

	returnErrorMock := func(context *gin.Context, statusCode int, err error) {
		assert.Fail(t, "An error happened, when is not expected error: "+err.Error())
	}

	NotificationHandler(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		getQueryFromParamsMock,
		returnJSONMock,
		returnErrorMock,
	)(nil)
}

func TestGetByID_UpdateFromError_OK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		assert.False(t, fullOrder)
		assert.True(t, payed)

		pay_id := "b10763c4-9f63-42d8-83ca-40824734ec6d"
		result := &modelDB.Order{
			ID:         "267511c9-324a-41dd-8711-c7d72e179522",
			Status:     "ERROR",
			Amount:     1000,
			Currency:   "BRL",
			CustomerID: "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			PaymentID:  &pay_id,
		}
		return result, 200, nil
	}

	getQueryFromParamsMock := func(context *gin.Context, params string) string {
		paramsMock := map[string]string{
			"order_id":    "267511c9-324a-41dd-8711-c7d72e179522",
			"status":      "APPROVED",
			"customer_id": "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			"payment_id":  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}

		param, ok := paramsMock[params]
		assert.True(t, ok, "Provided 'params' is invalid")
		assert.NotEmpty(t, param)

		return param
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		status, ok := changes["status"]
		assert.True(t, ok, "'status' need to be provided")
		assert.Equal(t, "APPROVED", status)
		_, ok = changes["payment_id"]
		assert.False(t, ok, "'payment_id' is not to be provided")
		return nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 200, statusCode)
		assert.Nil(t, obj)
	}

	returnErrorMock := func(context *gin.Context, statusCode int, err error) {
		assert.Fail(t, "An error happened, when is not expected error: "+err.Error())
	}

	NotificationHandler(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		getQueryFromParamsMock,
		returnJSONMock,
		returnErrorMock,
	)(nil)
}

func TestGetByID_PaymentIDNotMatch_NOK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		assert.False(t, fullOrder)
		assert.True(t, payed)

		pay_id := "b10763c4-9f63-42d8-83ca-40824734ec6d"
		result := &modelDB.Order{
			ID:         "267511c9-324a-41dd-8711-c7d72e179522",
			Status:     "ERROR",
			Amount:     1000,
			Currency:   "BRL",
			CustomerID: "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			PaymentID:  &pay_id,
		}
		return result, 200, nil
	}

	getQueryFromParamsMock := func(context *gin.Context, params string) string {
		paramsMock := map[string]string{
			"order_id":    "267511c9-324a-41dd-8711-c7d72e179522",
			"status":      "APPROVED",
			"customer_id": "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			"payment_id":  "fb3e8ebd-f09b-49c1-a8a7-7667737cdbec",
		}

		param, ok := paramsMock[params]
		assert.True(t, ok, "Provided 'params' is invalid")
		assert.NotEmpty(t, param)

		return param
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		assert.Fail(t, "An error happened, 'updateStatus' is not to be called")
		return nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Fail(t, "An error happened, 'returnJSON' is not to be called")
	}

	returnErrorMock := func(context *gin.Context, statusCode int, err error) {
		assert.Equal(t, 400, statusCode)
		assert.Error(t, err)
		assert.Equal(t, "\"payment_id\" does not match", err.Error())
	}

	NotificationHandler(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		getQueryFromParamsMock,
		returnJSONMock,
		returnErrorMock,
	)(nil)
}
