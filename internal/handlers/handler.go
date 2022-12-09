package handlers

import (
	"net_up_test_task/internal/models"

	"github.com/gin-gonic/gin"
)

// service Интерфейс к слою services, отвечающему за хранение данных в кэше
type service interface {
	Save(ip string)
	Get() []models.UserDTO
}

// HandlerConfig Конфигурация к обработчику
type HandlerConfig struct {
	Router	*gin.Engine
	Service	service
}

// Handler Обработчик запросов
type Handler struct {
	service	service
}

// RegisterHandler Регистратор обработчика
func RegisterHandler(c *HandlerConfig) {
	handler := &Handler{
		service: c.Service,
	}

	g := c.Router.Group("v1") // Версия API

	g.GET("/user/ping", handler.Ping)		// Пинг (для записи подключения)
	g.GET("/admin/users", handler.GetUsers)	// Юзерс (для получения активных соединений)
}
