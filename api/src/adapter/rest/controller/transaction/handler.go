package transactioncontroller

import (
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/payment/providers"
	"ravxcheckout/src/adapter/rest/helper"
	orderhelper "ravxcheckout/src/adapter/rest/helper/order"
	"ravxcheckout/src/internal/model"

	"github.com/gin-gonic/gin"
)

func Create(
	getOrderById order.GetByIDFn,
	getExistOrder orderhelper.GetExistOrderFn,
	updateStatus order.PatchFn,
	executePayment providers.ExecutePaymentFn,
	getObjectFromPostRequest helper.GetObjectFromPostRequestFn,
	returnJSON helper.ReturnJSONFn,
) gin.HandlerFunc {

	return func(context *gin.Context) {
		var paymentRequest model.PaymentRequest
		paymentResponse := &model.PaymentResponse{}

		err := getObjectFromPostRequest(context, &paymentRequest)
		if err != nil {
			returnJSON(context, 400, &model.PaymentResponse{
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
			})
			return
		}

		paymentResponse, err = executePayment(&paymentRequest, getOrderById, getExistOrder)
		if err != nil {
			returnJSON(context, 500, &model.PaymentResponse{
				Error: &model.Error{
					Message:    err.Error(),
					StatusCode: 500,
				},
			})
			return
		}

		statusCode := 201
		if paymentResponse.Error != nil {
			statusCode = paymentResponse.Error.StatusCode
		} else {
			changeStatus := map[string]interface{}{
				"status":     paymentResponse.Status,
				"payment_id": paymentResponse.PaymentID,
			}

			err = updateStatus(paymentResponse.OrderID, changeStatus)
			if err != nil {
				returnJSON(context, 500, &model.PaymentResponse{
					Error: &model.Error{
						Message:    err.Error(),
						StatusCode: 500,
					},
				})
				return
			}
		}

		returnJSON(context, statusCode, paymentResponse)
	}
}
