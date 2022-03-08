package pixcontroller

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
		// not used because GetExistOrder is mocked
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

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		return nil
	}

	executePixPaymentMock := func(
		paymentRequest *model.PixPaymentRequest,
		getOrderById order.GetByIDFn,
		getExistOrder helper.GetExistOrderFn,
	) (*model.PixPaymentResponse, error) {
		// success return
		paymentResponse := &model.PixPaymentResponse{
			PaymentID:   "1a938e0d-26ab-4ac2-b263-ab46a99e4356",
			Status:      "WAITING",
			Description: "QR Code gerado com sucesso e aguardando o pagamento.",
			AdditionalData: &model.AdditionalDataPix{
				TransactionID:        "8289874875871543653292342",
				QrCode:               "00020101021226740014br.gov.bcb.pix210812345678220412342308123456782420001122334455 667788995204000053039865406123.455802BR5913FULANO DE TAL6008BRASILIA62190515RP12345678- 201980720014br.gov.bcb.pix2550bx.com.br/spi/U0VHUkVET1RPVEFMTUVOVEVBTEVBVE9SSU8=63 0434D1",
				CreationDateQrCode:   "2021-02-14T17:50:12Z",
				ExpirationDateQrCode: "2021-02-14T17:51:52Z",
				PSPCode:              "033",
			},
		}
		return paymentResponse, nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 201, statusCode)
		paymentResponse, ok := obj.(*model.PixPaymentResponse)
		assert.Equal(t, true, ok)
		assert.Nil(t, paymentResponse.Error)
	}

	Create(
		getOrderByIdMock,
		getExistOrderMock,
		executePixPaymentMock,
		getObjectFromPostRequestMock,
		returnJSONMock,
	)(nil)
}

func TestCreate400_NOK(t *testing.T) {

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
		return nil, 200, nil
	}

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		return errors.New("Parsing Error")
	}

	executePixPaymentMock := func(
		paymentRequest *model.PixPaymentRequest,
		getOrderById order.GetByIDFn,
		getExistOrder helper.GetExistOrderFn,
	) (*model.PixPaymentResponse, error) {
		return nil, nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 400, statusCode)
		paymentResponse, ok := obj.(*model.PixPaymentResponse)
		assert.Equal(t, true, ok)
		assert.NotNil(t, paymentResponse.Error)
	}

	Create(
		getOrderByIdMock,
		getExistOrderMock,
		executePixPaymentMock,
		getObjectFromPostRequestMock,
		returnJSONMock,
	)(nil)
}

func TestCreateGetnetError_NOK(t *testing.T) {

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
		return nil, 200, nil
	}

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		return nil
	}

	executePixPaymentMock := func(
		paymentRequest *model.PixPaymentRequest,
		getOrderById order.GetByIDFn,
		getExistOrder helper.GetExistOrderFn,
	) (*model.PixPaymentResponse, error) {
		// fail return
		paymentResponse := &model.PixPaymentResponse{
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
		paymentResponse, ok := obj.(*model.PixPaymentResponse)
		assert.Equal(t, true, ok)
		assert.NotNil(t, paymentResponse.Error)
	}

	Create(
		getOrderByIdMock,
		getExistOrderMock,
		executePixPaymentMock,
		getObjectFromPostRequestMock,
		returnJSONMock,
	)(nil)
}
