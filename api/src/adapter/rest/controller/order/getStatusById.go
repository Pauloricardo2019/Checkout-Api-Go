package ordercontroller

import (
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/rest/helper"
	orderhelper "ravxcheckout/src/adapter/rest/helper/order"
	"ravxcheckout/src/internal/model/dto"

	"github.com/gin-gonic/gin"
)

func GetStatusByID(
	getOrderById order.GetByIDFn,
	getExistOrder orderhelper.GetExistOrderFn,
	getQueryFromParams helper.GetQueryFromParamsFn,
	returnJSON helper.ReturnJSONFn,
	returnError helper.ReturnErrorFn,
) gin.HandlerFunc {

	return func(context *gin.Context) {

		orderID := getQueryFromParams(context, "id")

		order, status, err := getExistOrder(orderID, false, true, getOrderById)
		if err != nil {
			returnError(context, status, err)
			return
		}

		returnJSON(context, 200, &dto.OrderDTO{
			Status: order.Status,
		})
	}

}
