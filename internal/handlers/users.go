package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsersResponse Структура ответа (по одному подключению)
type GetUsersResponse struct {
	IP				string	`json:"ip_address"`
	FirstRequest	string	`json:"since"`
}

// GetUsers Метод обработчика отдает список активных соединений
func (h *Handler) GetUsers(ctx *gin.Context) {
	// Получаем из кэша список соединений
	clients := h.service.Get()
	
	// Создаем массив ответа
	response := make([]GetUsersResponse, len(clients))

	// Маппим данные в ответ
	for k := 0; k < len (clients); k++ {
		response[k].IP = clients[k].IP
		response[k].FirstRequest = clients[k].FirstRequest.Format("2006-01-02 15:04:05")
	}

	// Отдаем ответ
	ctx.JSON(http.StatusOK, response)
}
