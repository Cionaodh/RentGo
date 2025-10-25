# rent_go:
# 	docker-compose -f docker-compose.yml up -d --force-recreate
rent_go:
	docker compose up --force-recreate --build -d

migrate:
	migrate create -ext sql -dir ./migrations -seq initial
