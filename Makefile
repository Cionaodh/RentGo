rent_go:
	docker-compose -f docker-compose.yml up -d --force-recreate

# migrate:
# 	migrate create -ext sql -dir ./migrations -seq create_rentpoint_table
	
# Swagger:
# 	swag init -g ./internal/controller/http/router.go --output docs --parseDependency