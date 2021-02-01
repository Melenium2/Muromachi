package main

import (
	"os"
)

// Даталоадер, агрегация постоянных одинаковых запросов
// https://gqlgen.com/reference/dataloaders/
//
// Сделать авторизацию (обсудить как удобнее) 	https://graphql.org/learn/authorization/
//												https://gqlgen.com/recipes/authentication/
// Про авторизацию
// https://gist.github.com/zmts/802dc9c3510d79fd40f9dc38a12bccfc
// Server - Client
// Есть ендпоинт. на него шлем реквест с типом application/x-www-form-urlencoded
// в хедере которого будет клиент id и клиент секрет как base64
// -> Первоначально отправляем с типом authorization_code
// -> Получаем body {
//   "access_token": "NgCXRK...MzYjw",
//   "token_type": "Bearer",
//   "expires_in": 3600,
//   "refresh_token": "NgAagA...Um_SHo"
//}
// -> refresh token сохраняем в куки дополнительно
// -> Чтоб рефрешнуть отправляем опять aспец хедер и тип refresh_token
// -> Так и делаем авторизацию
// -> Подробности почитать в гите
//		 https://gist.github.com/zmts/802dc9c3510d79fd40f9dc38a12bccfc

// Server - Server
// Есть ендпоинт. на него шлем реквест с типом application/x-www-form-urlencoded
// в хедере которого будет клиент id и клиент секрет как base64
// -> в ответ прилитает ответ {
//   "access_token": "NgCXRKc...MzYjw",
//   "token_type": "bearer",
//   "expires_in": 3600,
//}

// Показывать или нет ендпоинты https://gqlgen.com/reference/introspection/

// Сделать тулзу позволит лимитировать реквесты
// Создать отдельно сервер с логикой, сначала посмотреть есть ли его смысл
// Посмотреть как можно мокать fasthttp context в auth_test.go


const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//cfg := config.New("./config/dev.yml")
	//cfg.Database.Schema = "./config/schema.sql"
	//conn, err := store.EstablishConnection(cfg.Database)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//resolver := &graph.Resolver{
	//	Tables: store.New(conn),
	//}
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
