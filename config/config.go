package config

type Config struct {
	Token         string `json:"token"`          // Токен для Telegram бота
	RedisAddr     string `json:"redis_addr"`     // Адрес Redis
	RedisPassword string `json:"redis_password"` // Пароль для Redis
	RedisDB       int    `json:"redis_db"`       // Номер базы данных Redis
}

// NewConfig создаёт и возвращает новый экземпляр конфигурации.
func NewConfig() (*Config, error) {
	cfg := &Config{
		Token:         "8159719184:AAF-uQXUxzPXFjTS0A8-HR--FkEbsIvqtRg",
		RedisAddr:     "localhost:6379", // Локальный адрес Redis
		RedisPassword: "",               // Без пароля (если требуется, укажите здесь пароль)
		RedisDB:       0,                // Используем базу данных Redis с индексом 0
	}
	return cfg, nil
}
