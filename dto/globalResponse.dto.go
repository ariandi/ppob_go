package dto

import "github.com/gin-gonic/gin"

func ErrorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func ErrorResponseString(err string) gin.H {
	return gin.H{"message": err}
}

type ResponseDefault struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
