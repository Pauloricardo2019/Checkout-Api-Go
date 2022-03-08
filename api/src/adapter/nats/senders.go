package nats

import (
	"encoding/json"
	"ravxcheckout/src/internal/model"
	"time"
)

func SendPaymentResponse(response *model.PaymentResponse) error {
	conn, err := getNatsConn()
	if err != nil {
		panic(err.Error())
	}

	byteArray, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = conn.Publish(PaymentResponseQueue, byteArray)
	if err != nil {
		return err
	}

	return nil
}

func SendPaymentRequest(paymentRequest *model.PaymentRequest) error {
	conn, err := getNatsConn()
	if err != nil {
		panic(err.Error())
	}

	byteArray, err := json.Marshal(paymentRequest)
	if err != nil {
		return err
	}

	response, err := conn.Request(PaymentRequestQueue, byteArray, 30*time.Second)
	if err != nil {
		return err
	}

	paymentResponse := &model.PaymentResponse{}
	err = json.Unmarshal(response.Data, paymentResponse)
	if err != nil {
		return err
	}

	return nil
}
