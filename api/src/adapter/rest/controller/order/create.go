package ordercontroller

import (
	"ravxcheckout/src/adapter/config"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/rest/helper"
	model "ravxcheckout/src/internal/model/db"
	"ravxcheckout/src/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Create(
	persistOrder order.CreateFn,
	getObjectFromPostRequest helper.GetObjectFromPostRequestFn,
	getConfig config.GetConfigFn,
	returnJSON helper.ReturnJSONFn,
	returnError helper.ReturnErrorFn,
) gin.HandlerFunc {

	return func(context *gin.Context) {
		var order model.Order

		err := getObjectFromPostRequest(context, &order)
		if err != nil {
			returnError(context, 400, err)
			return
		}

		orderUUID := uuid.New().String()
		customerUUID := uuid.New().String()
		addressUUID := uuid.New().String()

		order.ID = orderUUID
		order.Customer.ID = customerUUID
		order.Customer.Address.ID = addressUUID

		err = persistOrder(&order)
		if err != nil {
			returnError(context, 500, err)
			return
		}

		ctg := getConfig()
		urlString := ctg.CheckoutURL + "?id=" + order.ID

		returnJSON(context, 200, &dto.OrderDTO{
			RedirectUrl: urlString,
		})
	}

}
