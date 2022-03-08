package notificationcontroller

import (
	"errors"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/rest/helper"
	orderhelper "ravxcheckout/src/adapter/rest/helper/order"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NotificationHandler(
	getOrderById order.GetByIDFn,
	getExistOrder orderhelper.GetExistOrderFn,
	updateStatus order.PatchFn,
	getQueryFromParams helper.GetQueryFromParamsFn,
	returnJSON helper.ReturnJSONFn,
	returnError helper.ReturnErrorFn,
) gin.HandlerFunc {

	return func(context *gin.Context) {

		orderID := getQueryFromParams(context, "order_id")

		status := getQueryFromParams(context, "status")
		if status == "" {
			returnError(context, 400, errors.New("\"status\" is not allowed to be empty"))
			return
		}

		customerID := getQueryFromParams(context, "customer_id")
		if customerID == "" {
			returnError(context, 400, errors.New("\"customer_id\" is not allowed to be empty"))
			return
		}

		paymentID := getQueryFromParams(context, "payment_id")
		if paymentID == "" {
			returnError(context, 400, errors.New("\"payment_id\" is not allowed to be empty"))
			return
		}

		_, err := uuid.Parse(paymentID)
		if err != nil {
			returnError(context, 400, errors.New("\"payment_id\" error: "+err.Error()))
			return
		}

		order, statusCode, err := getExistOrder(orderID, false, true, getOrderById)
		if err != nil {
			returnError(context, statusCode, err)
			return
		}

		if order.CustomerID != customerID {
			returnError(context, 400, errors.New("\"customer_id\" does not match"))
		}

		if order.Status == "APPROVED" || order.Status == status {
			returnJSON(context, 202, nil)
			return
		}

		changeStatus := map[string]interface{}{
			"status": status,
		}

		if order.PaymentID != nil {
			_, err := uuid.Parse(*order.PaymentID)
			if err != nil {
				changeStatus["payment_id"] = paymentID
			}

			if paymentID != *order.PaymentID {
				returnError(context, 400, errors.New("\"payment_id\" does not match"))
				return
			}

		} else {
			changeStatus["payment_id"] = paymentID
		}

		err = updateStatus(order.ID, changeStatus)
		if err != nil {
			returnError(context, 500, err)
			return
		}

		// need to send this status to Client
		returnJSON(context, 200, nil)

	}
}
