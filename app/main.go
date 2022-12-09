package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net_up_test_task/internal/handlers"
	"net_up_test_task/internal/services"

	"github.com/gin-gonic/gin"
)

// Конфигурация
const (
	Format	= "2006-01-02 15:04:05"	// Формат времени для ответа (/admin/users)
	TOut 	= time.Minute * 30		// Таймаут
	Host 	= "0.0.0.0"				// Адрес
	Port 	= "8080"				// Порт
)

func main() {
	// Инициализируем роутер
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Инициализируем сервис
	service := services.NewService(&services.ServiceConfig{
		Timeout:	TOut,
		TimeFormat:	Format,
	})

	// Регистрируем обработчик
	handlers.RegisterHandler(&handlers.HandlerConfig{
		Router:		router,
		Service:	service,
	})

	// Запускаем цикл в слое services
	channel := service.Run()
	log.Print("service started")

	// Инициализация сервера
	server := &http.Server{
		Addr:		fmt.Sprintf("%s:%s", Host, Port),
		Handler:	router,
	}

	// Запускаем сервер
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to initialize server: %v\n", err)
		}
	}()
	log.Printf("server listening on port %s", Port)

	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Отключаем сервис
	log.Print("shutting down service cycle...")
	channel <- struct{}{}

	// Отключаем сервер
	log.Print("shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v\n", err)
	}

	<-ctx.Done()
	log.Print("successfully")
}
