package providers

import (
	"ravxcheckout/src/adapter/database/sql/repository/order"
	helper "ravxcheckout/src/adapter/rest/helper/order"
	"ravxcheckout/src/internal/model"
)

type ExecutePaymentFn func(
	paymentRequest *model.PaymentRequest,
	getOrderById order.GetByIDFn,
	getExistOrder helper.GetExistOrderFn,
) (*model.PaymentResponse, error)

type ExecutePixPaymentFn func(
	paymentRequest *model.PixPaymentRequest,
	getOrderById order.GetByIDFn,
	getExistOrder helper.GetExistOrderFn,
) (*model.PixPaymentResponse, error)
