package getnet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ravxcheckout/src/adapter/config"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	helper "ravxcheckout/src/adapter/rest/helper/order"
	"ravxcheckout/src/internal/model"

	"github.com/google/uuid"
)

func ExecutePayment(
	paymentRequest *model.PaymentRequest,
	getOrderById order.GetByIDFn,
	getExistOrder helper.GetExistOrderFn,
) (*model.PaymentResponse, error) {

	if paymentRequest.Debit == nil {
		paymentResponse := &model.PaymentResponse{
			Error: &model.Error{
				Message:    "Bad Request",
				StatusCode: 400,
				Name:       "ValidationError",
				Details: []model.Details{
					{
						Status:            "DENIED",
						ErrorCode:         "GENERIC-400",
						Description:       "object is invalid",
						DescriptionDetail: "\"debit\" is not allowed to be empty",
					},
				},
			},
		}
		return paymentResponse, nil
	}

	order, status, err := getExistOrder(paymentRequest.OrderID, true, false, getOrderById)
	if err != nil {
		if status == 500 {
			return nil, err
		}

		if status == 404 {
			paymentResponse := &model.PaymentResponse{
				Error: &model.Error{
					Message:    "Not Found",
					StatusCode: 404,
					Name:       "ValidationError",
					Details: []model.Details{
						{
							Status:            "DENIED",
							ErrorCode:         "GENERIC-404",
							Description:       "Not found",
							DescriptionDetail: err.Error(),
						},
					},
				},
			}
			return paymentResponse, nil
		}

		paymentResponse := &model.PaymentResponse{
			Error: &model.Error{
				Message:    "Bad Request",
				StatusCode: 400,
				Name:       "ValidationError",
				Details: []model.Details{
					{
						Status:            "DENIED",
						ErrorCode:         "GENERIC-400",
						Description:       "object is invalid",
						DescriptionDetail: err.Error(),
					},
				},
			},
		}
		return paymentResponse, nil
	}

	paymentRequest.Order = order

	authResponse, err := authenticate()
	if err != nil {
		return nil, err
	}

	if authResponse.Error != nil {
		return nil, errors.New("authenticator error")
	}

	numberTokenResponse, err := tokenizeCard(authResponse, paymentRequest.Debit)
	if err != nil {
		return nil, err
	}

	if numberTokenResponse.Error != nil {
		paymentResponse := &model.PaymentResponse{
			Error: numberTokenResponse.Error,
		}
		return paymentResponse, nil
	}

	paymentResponse, err := paymentDebit(authResponse, numberTokenResponse, paymentRequest)
	if err != nil {
		return nil, err
	}

	return paymentResponse, nil
}

func getDeviceID(docNumber string) (string, error) {
	sessionID := fmt.Sprintf("%v-%v", docNumber, uuid.NewString())
	cfg := config.GetConfig()

	url := fmt.Sprintf("%v?org_id=%v&session_id=%v", cfg.DeviceURL, cfg.KeyDevice, sessionID)

	_, err := http.Get(url)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func ExecutePixPayment(
	paymentRequest *model.PixPaymentRequest,
	getOrderById order.GetByIDFn,
	getExistOrder helper.GetExistOrderFn,
) (*model.PixPaymentResponse, error) {

	order, status, err := getExistOrder(paymentRequest.OrderID, false, false, getOrderById)
	if err != nil {
		if status == 500 {
			return nil, err
		}

		if status == 404 {
			paymentResponse := &model.PixPaymentResponse{
				Error: &model.Error{
					Message:    "Not Found",
					StatusCode: 404,
					Name:       "ValidationError",
					Details: []model.Details{
						{
							Status:            "DENIED",
							ErrorCode:         "GENERIC-404",
							Description:       "Not found",
							DescriptionDetail: err.Error(),
						},
					},
				},
			}
			return paymentResponse, nil
		}

		paymentResponse := &model.PixPaymentResponse{
			Error: &model.Error{
				Message:    "Bad Request",
				StatusCode: 400,
				Name:       "ValidationError",
				Details: []model.Details{
					{
						Status:            "DENIED",
						ErrorCode:         "GENERIC-400",
						Description:       "object is invalid",
						DescriptionDetail: err.Error(),
					},
				},
			},
		}
		return paymentResponse, nil
	}

	paymentRequest.Order = order

	authResponse, err := authenticate()
	if err != nil {
		return nil, err
	}

	if authResponse.Error != nil {
		return nil, errors.New("Authenticator Error")
	}

	paymentResponse, err := paymentPix(authResponse, paymentRequest)
	if err != nil {
		return nil, err
	}

	return paymentResponse, nil
}

func authenticate() (*AuthResponse, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%v/auth/oauth/v2/token?scope=oob&grant_type=client_credentials", cfg.GetnetURL)

	authKeyStr := fmt.Sprintf("%v:%v", cfg.ClientID, cfg.ClientSecret)
	authKeyBase64 := base64.StdEncoding.EncodeToString([]byte(authKeyStr))

	auth := fmt.Sprintf("Basic %v", authKeyBase64)

	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Authorization", auth)

	client := &http.Client{}
	res, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	authResponse := &AuthResponse{}
	bodyDecoder := json.NewDecoder(res.Body)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if err := bodyDecoder.Decode(authResponse); err != nil {
			return nil, err
		}

		return authResponse, nil
	}

	errorResponse := &model.ErrorReduced{}
	if err := bodyDecoder.Decode(errorResponse); err != nil {
		return nil, err
	}
	authResponse.Error = errorResponse

	return authResponse, err
}

func paymentDebit(authResponse *AuthResponse, tokenizeCardResponse *TokenizeCardResponse, paymentRequest *model.PaymentRequest) (*model.PaymentResponse, error) {
	order := paymentRequest.Order
	customer := paymentRequest.Order.Customer

	deviceID, err := getDeviceID(customer.DocumentNumber)
	if err != nil {
		return nil, err
	}

	cfg := config.GetConfig()
	url := fmt.Sprintf("%v/v1/payments/debit", cfg.GetnetURL)

	debit := paymentRequest.Debit
	card := paymentRequest.Debit.Card
	addressOrder := customer.Address
	address := &model.Address{
		Street:     addressOrder.Street,
		Number:     addressOrder.Number,
		Complement: addressOrder.Complement,
		District:   addressOrder.District,
		City:       addressOrder.City,
		State:      addressOrder.State,
		Country:    addressOrder.Country,
		PostalCode: addressOrder.PostalCode,
	}

	payload := &Transaction{
		SellerID: cfg.SellerID,
		Amount:   order.Amount,
		Currency: order.Currency,
		Order: &model.Order{
			OrderID:     order.ID,
			ProductType: order.ProductType,
		},
		Customer: &model.Customer{
			CustomerID:     customer.ID,
			FirstName:      customer.FirstName,
			LastName:       customer.LastName,
			Name:           customer.Name,
			Email:          customer.Email,
			DocumentType:   customer.DocumentType,
			DocumentNumber: customer.DocumentNumber,
			PhoneNumber:    customer.PhoneNumber,
			Address:        address,
		},
		Device: &Device{
			DeviceID: deviceID,
		},
		Shippings: []model.Shippings{
			{
				FirstName:   customer.FirstName,
				Name:        customer.Name,
				Email:       customer.Email,
				PhoneNumber: customer.PhoneNumber,
				Address:     address,
			},
		},
		Debit: &Debit{
			CardHolderMobile: debit.CardHolderMobile,
			SoftDescription:  debit.SoftDescription,
			Authenticated:    debit.Authenticated,
			Card: &Card{
				NumberToken:     tokenizeCardResponse.NumberToken,
				CardHolderName:  card.CardHolderName,
				ExpirationMonth: card.ExpirationMonth,
				ExpirationYear:  card.ExpirationYear,
			},
		},
	}

	byteArrayPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteArrayPayload))
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("%v %v", authResponse.TokenType, authResponse.AccessToken))

	client := &http.Client{}
	res, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	paymentResponse := &model.PaymentResponse{}
	bodyDecoder := json.NewDecoder(res.Body)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if err := bodyDecoder.Decode(paymentResponse); err != nil {
			return nil, err
		}
		return paymentResponse, err
	}
	errorResponse := &model.Error{}
	if err := bodyDecoder.Decode(errorResponse); err != nil {
		return nil, err
	}
	paymentResponse.Error = errorResponse
	return paymentResponse, err

}

func tokenizeCard(authResponse *AuthResponse, debit *model.Debit) (*TokenizeCardResponse, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%v/v1/tokens/card", cfg.GetnetURL)

	payload := &CardTokenization{
		CardNumber: debit.Card.Number,
	}

	byteArrayPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteArrayPayload))
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Authorization", fmt.Sprintf("%v %v", authResponse.TokenType, authResponse.AccessToken))

	client := &http.Client{}
	res, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	numberTokenResponse := &TokenizeCardResponse{}
	bodyDecoder := json.NewDecoder(res.Body)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if err := bodyDecoder.Decode(numberTokenResponse); err != nil {
			return nil, err
		}
		return numberTokenResponse, nil
	}

	errorResponse := &model.Error{}
	if err := bodyDecoder.Decode(errorResponse); err != nil {
		return nil, err
	}
	numberTokenResponse.Error = errorResponse
	return numberTokenResponse, err

}

func paymentPix(authResponse *AuthResponse, paymentRequest *model.PixPaymentRequest) (*model.PixPaymentResponse, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%v/v1/payments/qrcode/pix", cfg.GetnetURL)

	order := paymentRequest.Order
	payload := &PixTransaction{
		Amount:     order.Amount,
		Currency:   order.Currency,
		OrderID:    order.ID,
		CustomerId: order.CustomerID,
	}

	byteArrayPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteArrayPayload))
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("%v %v", authResponse.TokenType, authResponse.AccessToken))

	client := &http.Client{}
	res, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	paymentResponse := &model.PixPaymentResponse{}
	bodyDecoder := json.NewDecoder(res.Body)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if err := bodyDecoder.Decode(paymentResponse); err != nil {
			return nil, err
		}
		return paymentResponse, err
	}
	errorResponse := &model.Error{}
	if err := bodyDecoder.Decode(errorResponse); err != nil {
		return nil, err
	}
	paymentResponse.Error = errorResponse
	return paymentResponse, err
}
