package models

import "time"

// User Структура хранения данных для одного активного соединения
type User struct {
	FirstRequest	time.Time	// Время первого запроса
	LastRequest		time.Time	// Время последнего запроса
	IP				string		// IP-адрес клиента
}
