package helper

import (
	modelDB "ravxcheckout/src/internal/model/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExistOrder_Full_OK(t *testing.T) {

	orderIDMock := "267511c9-324a-41dd-8711-c7d72e179522"
	fullOrderMock := true
	payedMock := true

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		assert.False(t, notFull)
		assert.True(t, payed)

		pay_id := "b10763c4-9f63-42d8-83ca-40824734ec6d"
		result := &modelDB.Order{
			ID:         "267511c9-324a-41dd-8711-c7d72e179522",
			Status:     "APPROVED",
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
			PaymentID: &pay_id,
		}
		return result, nil
	}

	order, statusCode, err := GetExistOrder(
		orderIDMock,
		fullOrderMock,
		payedMock,
		getOrderByIdMock,
	)

	assert.NotNil(t, order)
	assert.Equal(t, "267511c9-324a-41dd-8711-c7d72e179522", order.ID)
	assert.NotNil(t, order.Customer)
	assert.NotNil(t, order.Customer.Address)
	if order.Status == "APPROVED" {
		assert.NotNil(t, order.PaymentID)
	}
	assert.Equal(t, 200, statusCode)
	assert.Nil(t, err)

}

func TestGetExistOrder_NotFull_OK(t *testing.T) {

	orderIDMock := "267511c9-324a-41dd-8711-c7d72e179522"
	fullOrderMock := false
	payedMock := true

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		assert.True(t, notFull)
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
		return result, nil
	}

	order, statusCode, err := GetExistOrder(
		orderIDMock,
		fullOrderMock,
		payedMock,
		getOrderByIdMock,
	)

	assert.NotNil(t, order)
	assert.Equal(t, "267511c9-324a-41dd-8711-c7d72e179522", order.ID)
	assert.Nil(t, order.Customer)
	if order.Status == "APPROVED" {
		assert.NotNil(t, order.PaymentID)
	}
	assert.Equal(t, 200, statusCode)
	assert.Nil(t, err)

}

func TestGetExistOrder_FullNotPayed_OK(t *testing.T) {

	orderIDMock := "267511c9-324a-41dd-8711-c7d72e179522"
	fullOrderMock := true
	payedMock := false

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		assert.False(t, notFull)
		assert.False(t, payed)

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
		}
		return result, nil
	}

	order, statusCode, err := GetExistOrder(
		orderIDMock,
		fullOrderMock,
		payedMock,
		getOrderByIdMock,
	)

	assert.NotNil(t, order)
	assert.Equal(t, "267511c9-324a-41dd-8711-c7d72e179522", order.ID)
	assert.NotEqual(t, "APPROVED", order.Status)
	assert.NotNil(t, order.Customer)
	assert.NotNil(t, order.Customer.Address)
	assert.Equal(t, 200, statusCode)
	assert.Nil(t, err)

}

func TestGetExistOrder_Full_NOK(t *testing.T) {

	orderIDMock := "1234" // is not a UUID
	fullOrderMock := true
	payedMock := true

	getOrderByIdMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		assert.Fail(t, "Is not to execute this function")
		return nil, nil
	}

	order, statusCode, err := GetExistOrder(
		orderIDMock,
		fullOrderMock,
		payedMock,
		getOrderByIdMock,
	)

	assert.Nil(t, order)
	assert.Equal(t, 400, statusCode)
	assert.Error(t, err)
	assert.Equal(t, "\"order_id\" error: invalid UUID length: 4", err.Error())

}
