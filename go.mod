module Muromachi

go 1.15

replace github.com/99designs/gqlgen v0.13.0 => github.com/arsmn/gqlgen v0.13.2

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofiber/fiber/v2 v2.3.3
	github.com/jackc/pgconn v1.8.0
	github.com/jackc/pgproto3/v2 v2.0.6
	github.com/jackc/pgtype v1.6.2
	github.com/jackc/pgx/v4 v4.10.1
	github.com/stretchr/testify v1.5.1
	github.com/valyala/fasthttp v1.18.0
	github.com/vektah/gqlparser/v2 v2.1.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/yaml.v2 v2.2.4
)
