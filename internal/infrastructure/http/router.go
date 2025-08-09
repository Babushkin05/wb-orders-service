package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler Handler) {
	s := r.Group("/order")
	{
		s.GET("", handler.GetOrder)
	}
}
