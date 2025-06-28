RentGo:
	docker-compose build
	docker-compose up -d

Swagger:
	swag init -g ./internal/controller/http/router.go --output docs --parseDependency