package v1

import (
	"ravxcheckout/src/adapter/config"
	"ravxcheckout/src/adapter/database/sql/repository/order"
	"ravxcheckout/src/adapter/payment/providers/getnet"
	notificationcontroller "ravxcheckout/src/adapter/rest/controller/notification"
	ordercontroller "ravxcheckout/src/adapter/rest/controller/order"
	pixcontroller "ravxcheckout/src/adapter/rest/controller/pix"
	transactioncontroller "ravxcheckout/src/adapter/rest/controller/transaction"
	"ravxcheckout/src/adapter/rest/helper"
	orderhelper "ravxcheckout/src/adapter/rest/helper/order"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(router *gin.RouterGroup) {

	orderRouter := router.Group("/order")
	{
		orderRouter.POST("",
			ordercontroller.Create(
				order.Create,
				helper.GetObjectFromPostRequest,
				config.GetConfig,
				helper.ReturnJSON,
				helper.ReturnError,
			))
		orderRouter.GET("",
			ordercontroller.GetByID(
				order.GetByID,
				orderhelper.GetExistOrder,
				helper.GetQueryFromParams,
				config.GetConfig,
				helper.ReturnJSON,
				helper.ReturnError,
			))
		orderRouter.GET("/status",
			ordercontroller.GetStatusByID(
				order.GetByID,
				orderhelper.GetExistOrder,
				helper.GetQueryFromParams,
				helper.ReturnJSON,
				helper.ReturnError,
			))
	}
	transactionRouter := router.Group("/transaction")
	{
		transactionRouter.POST("",
			transactioncontroller.Create(
				order.GetByID,
				orderhelper.GetExistOrder,
				order.Patch,
				getnet.ExecutePayment,
				helper.GetObjectFromPostRequest,
				helper.ReturnJSON,
			))
	}
	pixRouter := router.Group("/pix")
	{
		pixRouter.POST("",
			pixcontroller.Create(
				order.GetByID,
				orderhelper.GetExistOrder,
				getnet.ExecutePixPayment,
				helper.GetObjectFromPostRequest,
				helper.ReturnJSON,
			))
	}
	notificationRouter := router.Group("/notification")
	{
		notificationRouter.GET("/debit",
			notificationcontroller.NotificationHandler(
				order.GetByID,
				orderhelper.GetExistOrder,
				order.Patch,
				helper.GetQueryFromParams,
				helper.ReturnJSON,
				helper.ReturnError,
			))
		notificationRouter.GET("/pix",
			notificationcontroller.NotificationHandler(
				order.GetByID,
				orderhelper.GetExistOrder,
				order.Patch,
				helper.GetQueryFromParams,
				helper.ReturnJSON,
				helper.ReturnError,
			))
	}

}
