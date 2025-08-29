package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoverJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				Fail(c, http.StatusInternalServerError, "internal server error")
			}
		}()
		c.Next()
	}
}

func NotFoundAsJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == http.StatusNotFound && !c.Writer.Written() {
			Fail(c, http.StatusNotFound, "not found")
		}
	}
}
