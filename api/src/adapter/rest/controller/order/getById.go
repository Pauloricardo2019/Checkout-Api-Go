package ordercontroller

import (
	"ravxcheckout/src/adapter/config"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/rest/helper"
	orderhelper "ravxcheckout/src/adapter/rest/helper/order"

	"github.com/gin-gonic/gin"
)

func GetByID(
	getOrderById order.GetByIDFn,
	getExistOrder orderhelper.GetExistOrderFn,
	getQueryFromParams helper.GetQueryFromParamsFn,
	getConfig config.GetConfigFn,
	returnJSON helper.ReturnJSONFn,
	returnError helper.ReturnErrorFn,
) gin.HandlerFunc {

	return func(context *gin.Context) {

		orderID := getQueryFromParams(context, "id")

		order, status, err := getExistOrder(orderID, true, true, getOrderById)
		if err != nil {
			returnError(context, status, err)
			return
		}

		urlString := getConfig().RedirectUrl
		order.RedirectUrl = &urlString

		returnJSON(context, 200, order)
	}

}
