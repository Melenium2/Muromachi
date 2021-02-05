package main

import (
	"Muromachi/config"
	"Muromachi/server"
	"log"
	"os"
	"os/signal"
)

// Даталоадер, агрегация постоянных одинаковых запросов
// https://gqlgen.com/reference/dataloaders/
//
// Показывать или нет ендпоинты https://gqlgen.com/reference/introspection/
//
// Сделать тулзу позволит лимитировать реквесты
//		расставить разные лимиты для разных категорий токенов
// Посмотреть как можно мокать fasthttp context в auth_test.go
// Дописать конфиг. Чтобы брать новые envs для redis и соль для авторизации
// TODO Изменить возвращающийся ошибки, сейчас какой то рандом
//		потом поправить тесты
// Возможно swagger для rest http
// Check db connection before server started
// Логирование всех запросов в бд

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cfg := config.New("./config/dev.yml")
	cfg.Database.Schema = "./config/schema.sql"

	serv := server.New(defaultPort, cfg)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		_ = <-c
		_ = serv.Shutdown()
	}()

	if err := serv.Listen(); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutdown")
}
