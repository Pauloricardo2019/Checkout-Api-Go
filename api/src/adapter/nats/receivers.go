package nats

import (
	"encoding/json"
	"fmt"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/monitoring"
	"ravxcheckout/src/adapter/payment/providers/getnet"
	helper "ravxcheckout/src/adapter/rest/helper/order"
	"ravxcheckout/src/internal/model"

	"github.com/nats-io/nats.go"
)

var subscriptions []*nats.Subscription

func init() {
	subscriptions = []*nats.Subscription{}
}

func receivePaymentRequest(msg *nats.Msg) {
	paymentRequest := &model.PaymentRequest{}
	paymentResponse := &model.PaymentResponse{}

	err := json.Unmarshal(msg.Data, paymentRequest)
	if err != nil {
		fmt.Printf("receiveUser: ERROR - %v\n", err.Error())
		monitoring.CaptureError(err)
		respondError(msg, err)
		return
	}

	paymentResponse, err = getnet.ExecutePayment(paymentRequest, order.GetByID, helper.GetExistOrder)
	if err != nil {
		monitoring.CaptureError(err)
		respondError(msg, err)
		return
	}

	respondMsg(msg, paymentResponse)
	return
}

func respondMsg(msg *nats.Msg, paymentResponse *model.PaymentResponse) {
	result, err := json.Marshal(paymentResponse)
	if err != nil {
		return
	}

	err = msg.Respond(result)
	if err != nil {
		monitoring.CaptureError(err)
	}
}

func respondError(msg *nats.Msg, err error) {
	paymentResponse := &model.PaymentResponse{Error: &model.Error{Message: err.Error()}}

	result, err := json.Marshal(paymentResponse)
	if err != nil {
		return
	}

	err = msg.Respond(result)
	if err != nil {
		monitoring.CaptureError(err)
	}
}

func StartSubscriptions() *nats.Conn {
	conn, err := getNatsConn()
	if err != nil {
		monitoring.CaptureError(err)
		panic(err.Error())
	}

	subscription, err := conn.Subscribe(PaymentRequestQueue, receivePaymentRequest)
	if err != nil {
		monitoring.CaptureError(err)
		panic(err.Error())
	}

	subscriptions = append(subscriptions, subscription)

	return conn
}

func StopSubscriptions() {
	for _, subscription := range subscriptions {
		err := subscription.Drain()
		if err != nil {
			monitoring.CaptureError(err)
			fmt.Printf("Error draining subscriptions. ERROR: %v", err.Error())
			continue
		}
	}
}
