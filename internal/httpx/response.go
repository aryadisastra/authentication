package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Result  bool        `json:"result"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Result: true, Data: data, Message: "Success Get Data"})
}
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Result: true, Data: data, Message: "Success Post Data"})
}
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(code, Envelope{Result: false, Data: gin.H{}, Message: msg})
}
