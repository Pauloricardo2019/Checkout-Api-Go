package ordercontroller

import (
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/internal/model"
	modelDB "ravxcheckout/src/internal/model/db"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetByID_OK(t *testing.T) {

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		result := &modelDB.Order{}
		return result, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		assert.True(t, fullOrder)
		assert.True(t, payed)

		result := &modelDB.Order{
			ID:         "267511c9-324a-41dd-8711-c7d72e179522",
			Status:     "WAITING_PAYMENT",
			Amount:     1000,
			Currency:   "BRL",
			CustomerID: "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
			Customer: &modelDB.Customer{
				ID:             "bc922b7f-94e0-4b47-b76a-6ede9bfb566d",
				FirstName:      "Eduardo",
				LastName:       "Andrade",
				Name:           "Eduardo Andrade",
				Email:          "eduardo@andrade.com",
				DocumentType:   "CPF",
				DocumentNumber: "12345678912",
				PhoneNumber:    "",
				AddressID:      "d52d0942-84a0-4648-88ff-f7681eed2b2f",
				Address: &modelDB.Address{
					ID:         "d52d0942-84a0-4648-88ff-f7681eed2b2f",
					Street:     "Rua João Rodrigues Leite",
					Number:     "196",
					Complement: "Bloco 2",
					District:   "Chácara Inglesa",
					City:       "São Paulo",
					State:      "SP",
					Country:    "Brazil",
					PostalCode: "05141160",
				},
			},
			// PaymentID:  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}
		return result, 200, nil
	}

	getQueryFromParamsMock := func(context *gin.Context, params string) string {
		return ""
	}

	getConfigMock := func() *model.Config {
		result := &model.Config{
			RedirectUrl: "test-url",
		}
		return result
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		order, ok := obj.(*modelDB.Order)
		assert.True(t, ok)
		assert.Equal(t, "test-url", *order.RedirectUrl)
	}

	returnErrorMock := func(context *gin.Context, statusCode int, err error) {
		assert.Fail(t, "An error happened, when is not expected error: "+err.Error())
	}

	GetByID(
		getOrderByIdMock,
		getExistOrderMock,
		getQueryFromParamsMock,
		getConfigMock,
		returnJSONMock,
		returnErrorMock,
	)(nil)
}
