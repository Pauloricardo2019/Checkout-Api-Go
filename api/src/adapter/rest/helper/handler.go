package helper

import (
	"ravxcheckout/src/internal/model/dto"

	"github.com/gin-gonic/gin"
)

type GetObjectFromPostRequestFn func(context *gin.Context, obj interface{}) error

func GetObjectFromPostRequest(context *gin.Context, obj interface{}) error {
	return context.ShouldBind(obj)
}

type GetQueryFromParamsFn func(context *gin.Context, params string) string

func GetQueryFromParams(context *gin.Context, params string) string {
	return context.Query(params)
}

type GetValueFromPathFn func(context *gin.Context, path string) string

func GetValueFromPath(context *gin.Context, path string) string {
	return context.Param(path)
}

type ReturnJSONFn func(context *gin.Context, statusCode int, obj interface{})

func ReturnJSON(context *gin.Context, statusCode int, obj interface{}) {
	context.JSON(statusCode, obj)
}

type ReturnErrorFn func(context *gin.Context, statusCode int, err error)

func ReturnError(context *gin.Context, statusCode int, err error) {
	context.JSON(statusCode, &dto.ErrorDTO{
		Message: err.Error(),
	})
}
