package http

import (
	"net/http"

	"github.com/Babushkin05/wb-orders-service/internal/application"
	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
	"github.com/Babushkin05/wb-orders-service/internal/shared/dto"
	"github.com/Babushkin05/wb-orders-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetOrder(c *gin.Context)
}

type handler struct {
	service application.OrdersService
}

func NewHandler(service application.OrdersService) Handler {
	return &handler{
		service: service,
	}
}

// @Summary Получить заказ по ID
// @Description Возвращает информацию о заказе по его OrderUID
// @Tags orders
// @Produce json
// @Param order_uid path string true "Order UID"
// @Success 200 {object} dto.Order "Успешный ответ"
// @Failure 404 {object} dto.ErrorResponse "Заказ не найден"
// @Router /orders/{order_uid} [get]
func (h *handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Infof("GetSubscription: getting subscription %s", id)

	order, err := h.service.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "subscription not found"})
		return
	}

	resp, err := model.MarshalOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	logger.Log.Infof("GetSubscription: found subscription %s", id)
	c.Data(http.StatusOK, "application/json", resp)
}
