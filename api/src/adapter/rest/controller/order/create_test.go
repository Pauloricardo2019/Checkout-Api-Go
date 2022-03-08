package ordercontroller

import (
	"ravxcheckout/src/internal/model"
	modelDB "ravxcheckout/src/internal/model/db"
	"ravxcheckout/src/internal/model/dto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate_OK(t *testing.T) {

	getObjectFromPostRequestMock := func(context *gin.Context, obj interface{}) error {
		order, ok := obj.(*modelDB.Order)
		assert.Equal(t, true, ok)

		*order = modelDB.Order{
			Status:   "WAITING_PAYMENT",
			Amount:   1000,
			Currency: "BRL",
			Customer: &modelDB.Customer{
				FirstName:      "Eduardo",
				LastName:       "Andrade",
				Name:           "Eduardo Andrade",
				Email:          "eduardo@andrade.com",
				DocumentType:   "CPF",
				DocumentNumber: "12345678912",
				PhoneNumber:    "",
				Address: &modelDB.Address{
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
		}
		return nil
	}

	getConfigMock := func() *model.Config {
		return &model.Config{
			CheckoutURL: "myURL",
		}
	}

	persistOrderMock := func(order *modelDB.Order) error {
		_, err := uuid.Parse(order.ID)
		assert.NoError(t, err)
		_, err = uuid.Parse(order.Customer.ID)
		assert.NoError(t, err)
		_, err = uuid.Parse(order.Customer.Address.ID)
		assert.NoError(t, err)

		order.ID = "267511c9-324a-41dd-8711-c7d72e179522"
		order.CustomerID = order.Customer.ID
		order.Customer.AddressID = order.Customer.Address.ID
		return nil
	}

	returnJSONMock := func(context *gin.Context, statusCode int, obj interface{}) {
		assert.Equal(t, 200, statusCode)
		orderResponse, ok := obj.(*dto.OrderDTO)
		assert.True(t, ok)
		assert.Equal(t, "myURL?id=267511c9-324a-41dd-8711-c7d72e179522", orderResponse.RedirectUrl)
	}

	returnErrorMock := func(context *gin.Context, statusCode int, err error) {
		assert.Fail(t, "An error happened, when is not expected error: "+err.Error())
	}

	Create(
		persistOrderMock,
		getObjectFromPostRequestMock,
		getConfigMock,
		returnJSONMock,
		returnErrorMock,
	)(nil)
}
