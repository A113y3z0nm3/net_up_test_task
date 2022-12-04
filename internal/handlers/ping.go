package handlers

import (
	"net/http"
	

	"github.com/gin-gonic/gin"
)

// Ping Метод обработчика записывает активные соединения в кэш
func (h *Handler) Ping(ctx *gin.Context) {
	// Получаем IP-адрес клиента
	ip := ctx.ClientIP()

	// Сохраняем в базу
	h.service.Save(ip)
	
	// Отдаем ответ
	ctx.JSON(http.StatusOK, gin.H{})
}
