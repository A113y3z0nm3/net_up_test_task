package services

import (
	"sync"
	"time"
	"net_up_test_task/internal/models"
)

// ServiceConfig Конфигурация для сервиса
type ServiceConfig struct {
	Timeout	time.Duration			// Таймаут (задается в main.go)
}

// Service Сервис, отвечающий за хранение и управление активными соединениями
type Service struct {
	timeout	time.Duration			// Таймаут (задается в ServiceConfig)
	mux		sync.RWMutex			// Мютекс (для конкуррентного доступа)
	clients	map[string]models.User	// Список активных соединений (клиентов)
}

// NewService Фабрика для сервиса
func NewService(c *ServiceConfig) *Service {
	return &Service{
		timeout:	c.Timeout,
		mux:		sync.RWMutex{},
		clients:	make(map[string]models.User),
	}
}

// Run Основная функция сервиса для мониторинга статусов подключений (продление, удаление)
func (s *Service) Run() chan struct{} {
	// Канал остановки (Для shutdown'а)
	doneChannel := make(chan struct{}, 1)
	// Тикер (Для проверок)
	ticker := time.NewTicker(500 * time.Millisecond)

	// Запускаем основной цикл Run
	go func(doneChannel chan struct{}, ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				// Проверяем подключения
				s.check()
			case <-doneChannel:
				// Останавливаем ticker
				ticker.Stop()

				return
			}
		}
	}(doneChannel, ticker)

	return doneChannel
}

// Check Проверяет активные подключения по таймауту
func (s *Service) check() {
	s.mux.RLock()

	// Итерируемся по списку активных подключений
	for ip, client := range s.clients {
		// Если прошло времени больше, чем заданный таймаут - удаляем клиента из кэша
		if time.Since(client.LastRequest) > s.timeout {
			s.mux.RUnlock()
			s.delete(ip)
			s.mux.RLock()
		}
	}

	s.mux.RUnlock()
}

// Save Сохраняет подключение, либо, если оно уже в списке активных, обновляет время последнего запроса
func (s *Service) Save(ip string) {
	// Проверяем адрес на наличие в кэше
	s.mux.RLock()
	client, ok := s.clients[ip]
	s.mux.RUnlock()

	if ok {
		// Если найден - обновляем дату последнего запроса
		client.LastRequest = time.Now()
	} else {
		// Если нет - записываем новое активное подключение
		client = models.User{
			FirstRequest: time.Now(),
			LastRequest: time.Now(),
			IP: ip,
		}
	}

	// Записываем в кэш
	s.mux.Lock()
	s.clients[ip] = client
	s.mux.Unlock()
}

// Delete Удаляет подключение из кэша
func (s *Service) delete(ip string) {
	// Удаляем по ключу (айпи)
	s.mux.Lock()
	delete(s.clients, ip)
	s.mux.Unlock()
}

// Get Отдает список активных соединений из кэша
func (s *Service) Get() []models.User {
	s.mux.RLock()
	
	// Создаем итоговый список
	result := make([]models.User, 0, len(s.clients))

	// Итерируемся по списку подключений
	for _, c := range s.clients {
		// Добавляем клиентов в итоговый список
		result = append(result, c)
	}

	s.mux.RUnlock()

	return result
}
