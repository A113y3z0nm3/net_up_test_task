package models

import "time"

// UserDTO Структура для перемещения данных между слоями для одного активного соединения
type UserDTO struct {
	FirstRequest	string	// Время первого запроса
	IP				string		// IP-адрес клиента
}

// User Структура для хранения данных в кэше для одного активного соединения
type User struct {
	FirstRequest	time.Time	// Время первого запроса
	LastRequest		time.Time	// Время последнего запроса
}
