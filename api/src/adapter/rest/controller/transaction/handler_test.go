package transactioncontroller

import (
	"errors"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	helper "ravxcheckout/src/adapter/rest/helper/order"
	"ravxcheckout/src/internal/model"
	modelDB "ravxcheckout/src/internal/model/db"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreate_OK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		return nil, 200, nil
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		assert.Equal(t, map[string]interface{}{
			"status":     "PENDING",
			"payment_id": "1234",
		}, changes)
		return nil
	}

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		return nil
	}

	executePaymentMock := func(
		paymentRequest *model.PaymentRequest,
		getOrderById order.GetByIDFn,
		getExistOrder helper.GetExistOrderFn,
	) (*model.PaymentResponse, error) {
		// success return
		paymentResponse := &model.PaymentResponse{
			PaymentID:   "1234",
			SellerID:    "12",
			RedirectUrl: "redirect_url",
			Status:      "PENDING",
			PostData: &model.PostData{
				IssuerPaymentID:            "1235",
				PayerAuthenticationRequest: "authenticated_code",
			},
		}
		return paymentResponse, nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 201, statusCode)
		paymentResponse, ok := obj.(*model.PaymentResponse)
		assert.Equal(t, true, ok)
		assert.Nil(t, paymentResponse.Error)
	}

	Create(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		executePaymentMock,
		getObjectFromPostRequestMock,
		returnJSONMock,
	)(nil)
}

func TestCreate400_NOK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		return nil, 200, nil
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		return nil
	}

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		return errors.New("Parsing Error")
	}

	executePaymentMock := func(
		paymentRequest *model.PaymentRequest,
		getOrderById order.GetByIDFn,
		getExistOrder helper.GetExistOrderFn,
	) (*model.PaymentResponse, error) {
		return nil, nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 400, statusCode)
		paymentResponse, ok := obj.(*model.PaymentResponse)
		assert.Equal(t, true, ok)
		assert.NotNil(t, paymentResponse.Error)
	}

	Create(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		executePaymentMock,
		getObjectFromPostRequestMock,
		returnJSONMock,
	)(nil)
}

func TestCreateGetnetError_NOK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		return nil, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		return nil, 200, nil
	}

	updateStatusMock := func(ID string, changes map[string]interface{}) error {
		return nil
	}

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		return nil
	}

	executePaymentMock := func(
		paymentRequest *model.PaymentRequest,
		getOrderById order.GetByIDFn,
		getExistOrder helper.GetExistOrderFn,
	) (*model.PaymentResponse, error) {
		// fail return
		paymentResponse := &model.PaymentResponse{
			Error: &model.Error{
				Message:    "message_error",
				Name:       "name_error",
				StatusCode: 400,
				Details:    nil,
			},
		}
		return paymentResponse, nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 400, statusCode)
		paymentResponse, ok := obj.(*model.PaymentResponse)
		assert.Equal(t, true, ok)
		assert.NotNil(t, paymentResponse.Error)
	}

	Create(
		getOrderByIdMock,
		getExistOrderMock,
		updateStatusMock,
		executePaymentMock,
		getObjectFromPostRequestMock,
		returnJSONMock,
	)(nil)
}
