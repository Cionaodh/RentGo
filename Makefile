RentGo:
	docker-compose -f docker-compose.yml up -d --force-recreate

# Swagger:
# 	swag init -g ./internal/controller/http/router.go --output docs --parseDependency