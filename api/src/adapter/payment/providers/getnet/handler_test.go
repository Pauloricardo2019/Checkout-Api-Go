package getnet

import (
	"encoding/json"
	"fmt"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/internal/model"
	modelDB "ravxcheckout/src/internal/model/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationExecutePayment(t *testing.T) {

	getByIDMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
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
			// PaymentID:  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}
		return result, 200, nil
	}

	request := &model.PaymentRequest{
		OrderID: "267511c9-324a-41dd-8711-c7d72e179522",
		Debit: &model.Debit{
			CardHolderMobile: "5551999887766",
			SoftDescription:  "LOJA*TESTE*COMPRA-123",
			Authenticated:    false,
			Card: &model.Card{
				Number:          "4012001037141112",
				CardHolderName:  "JOAO DA SILVA",
				SecurityCode:    "123",
				ExpirationMonth: "12",
				ExpirationYear:  "28",
			},
		},
	}

	paymentResponse, err := ExecutePayment(
		request,
		getByIDMock,
		getExistOrderMock,
	)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

	assert.NotEmpty(t, paymentResponse.PaymentID)
	assert.NotEmpty(t, paymentResponse.SellerID)

	_, err = json.Marshal(paymentResponse)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

}

func TestIntegrationInvalidCardExecutePayment(t *testing.T) {

	getByIDMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
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
			// PaymentID:  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}
		return result, 200, nil
	}

	request := &model.PaymentRequest{
		OrderID: "267511c9-324a-41dd-8711-c7d72e179522",
		Debit: &model.Debit{
			CardHolderMobile: "5551999887766", //5155901222270002
			SoftDescription:  "LOJA*TESTE*COMPRA-123",
			Authenticated:    false,
			Card: &model.Card{
				Number:          "4012001037141333",
				CardHolderName:  "JOAO DA SILVA",
				SecurityCode:    "123",
				ExpirationMonth: "12",
				ExpirationYear:  "28",
			},
		},
	}

	paymentResponse, err := ExecutePayment(
		request,
		getByIDMock,
		getExistOrderMock,
	)
	if err != nil {
		fmt.Printf("err1: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

	assert.Empty(t, paymentResponse.PaymentID)
	assert.NotEmpty(t, paymentResponse.Error)

	resError := paymentResponse.Error
	resErrorDetail := resError.Details[0]

	assert.Equal(t, 402, resError.StatusCode)
	assert.Equal(t, "GetnetTransactionError", resError.Name)
	assert.Equal(t, "PAYMENTS-077", resErrorDetail.ErrorCode)
	assert.Equal(t, "Numero de cartão inválido.", resErrorDetail.DescriptionDetail)
}

func TestIntegrationExpiredCardExecutePayment(t *testing.T) {

	getByIDMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
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
			// PaymentID:  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}
		return result, 200, nil
	}

	request := &model.PaymentRequest{
		OrderID: "267511c9-324a-41dd-8711-c7d72e179522",
		Debit: &model.Debit{
			Authenticated: false,
			Card: &model.Card{
				Number:          "4012001037141112",
				CardHolderName:  "JOAO DA SILVA",
				SecurityCode:    "123",
				ExpirationMonth: "12",
				ExpirationYear:  "19",
			},
		},
	}

	paymentResponse, err := ExecutePayment(
		request,
		getByIDMock,
		getExistOrderMock,
	)
	if err != nil {
		fmt.Printf("err1: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

	assert.Empty(t, paymentResponse.PaymentID)
	assert.NotEmpty(t, paymentResponse.Error)

	resError := paymentResponse.Error
	resErrorDetail := resError.Details[0]

	assert.Equal(t, "GetnetTransactionError", resError.Name)
	assert.Equal(t, "PAYMENTS-999", resErrorDetail.ErrorCode)
	assert.Equal(t, "ERROR - CGW000289-Card Expiration Day Invalid", resErrorDetail.DescriptionDetail)
}

func TestIntegrationEmptyRequiredFieldExecutePayment(t *testing.T) {

	getByIDMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
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
			// PaymentID:  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}
		return result, 200, nil
	}

	request := &model.PaymentRequest{
		OrderID: "267511c9-324a-41dd-8711-c7d72e179522",
		Debit: &model.Debit{
			Authenticated: false,
			Card: &model.Card{
				Number:         "4012001037141112",
				SecurityCode:   "123",
				ExpirationYear: "28",
			},
		},
	}

	paymentResponse, err := ExecutePayment(
		request,
		getByIDMock,
		getExistOrderMock,
	)
	if err != nil {
		fmt.Printf("err1: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

	assert.Empty(t, paymentResponse.PaymentID)
	assert.NotEmpty(t, paymentResponse.Error)
	assert.Equal(t, "\"expiration_month\" is not allowed to be empty", paymentResponse.Error.Details[0].DescriptionDetail)
}

func TestIntegrationExecutePixPayment(t *testing.T) {

	getByIDMock := func(ID string, notFull bool, payed bool) (*modelDB.Order, error) {
		// not used because GetExistOrder is mocked
		result := &modelDB.Order{}
		return result, nil
	}

	getExistOrderMock := func(
		orderID string,
		fullOrder bool,
		payed bool,
		getOrderById order.GetByIDFn,
	) (*modelDB.Order, int, error) {
		assert.False(t, fullOrder)
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
			// PaymentID:  "b10763c4-9f63-42d8-83ca-40824734ec6d",
		}
		return result, 200, nil
	}

	request := &model.PixPaymentRequest{
		OrderID: "267511c9-324a-41dd-8711-c7d72e179522",
	}

	paymentResponse, err := ExecutePixPayment(
		request,
		getByIDMock,
		getExistOrderMock,
	)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

	assert.NotEmpty(t, paymentResponse.PaymentID)
	assert.NotEmpty(t, paymentResponse.Status)
	assert.NotEmpty(t, paymentResponse.Description)
	assert.NotEmpty(t, paymentResponse.AdditionalData)
	assert.NotEmpty(t, paymentResponse.AdditionalData.TransactionID)
	assert.NotEmpty(t, paymentResponse.AdditionalData.QrCode)
	assert.NotEmpty(t, paymentResponse.AdditionalData.CreationDateQrCode)
	assert.NotEmpty(t, paymentResponse.AdditionalData.ExpirationDateQrCode)
	assert.NotEmpty(t, paymentResponse.AdditionalData.PSPCode)

	assert.Empty(t, paymentResponse.Error)

	_, err = json.Marshal(paymentResponse)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
		assert.Fail(t, err.Error())
		return
	}

}

// func TestIntegrationEmptyRequiredFieldExecutePixPayment(t *testing.T) {
// 	request := &model.PixPaymentRequest{
// 		Amount:     100,
// 		OrderID:    "DEV-1608748980",
// 		CustomerId: "string",
// 	}

// 	paymentResponse, err := ExecutePixPayment(request)
// 	if err != nil {
// 		fmt.Printf("err: %v", err.Error())
// 		assert.Fail(t, err.Error())
// 		return
// 	}

// 	assert.NotEmpty(t, paymentResponse.PaymentID)
// 	assert.Equal(t, "ERROR", paymentResponse.Status)
// 	assert.Empty(t, paymentResponse.Description)
// 	assert.Empty(t, paymentResponse.AdditionalData)

// 	assert.Empty(t, paymentResponse.Error)

// 	_, err = json.Marshal(paymentResponse)
// 	if err != nil {
// 		fmt.Printf("err: %v", err.Error())
// 		assert.Fail(t, err.Error())
// 		return
// 	}
// }

// func TestIntegrationBadRequestExecutePixPayment(t *testing.T) {
// 	request := &model.PixPaymentRequest{
// 		OrderID:    "DEV-1608748980",
// 		CustomerId: "string",
// 	}

// 	paymentResponse, err := ExecutePixPayment(request)
// 	if err != nil {
// 		fmt.Printf("err: %v", err.Error())
// 		assert.Fail(t, err.Error())
// 		return
// 	}

// 	assert.Empty(t, paymentResponse.PaymentID)
// 	assert.Empty(t, paymentResponse.Status)
// 	assert.Empty(t, paymentResponse.Description)
// 	assert.Empty(t, paymentResponse.AdditionalData)

// 	assert.Equal(t, "Bad Request", paymentResponse.Error.Message)
// 	assert.Equal(t, "ValidationError", paymentResponse.Error.Name)
// 	assert.Equal(t, 400, paymentResponse.Error.StatusCode)
// 	assert.Equal(t, "DENIED", paymentResponse.Error.Details[0].Status)
// 	assert.Equal(t, "GENERIC-400", paymentResponse.Error.Details[0].ErrorCode)
// 	assert.Equal(t, "Bad Request", paymentResponse.Error.Details[0].Description)
// 	assert.Equal(t, "Validation Error", paymentResponse.Error.Details[0].DescriptionDetail)
// 	assert.Empty(t, paymentResponse.Error.Details[0].Antifraud)

// 	_, err = json.Marshal(paymentResponse)
// 	if err != nil {
// 		fmt.Printf("err: %v", err.Error())
// 		assert.Fail(t, err.Error())
// 		return
// 	}
// }
