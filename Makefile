rent_go:
	docker-compose -f docker-compose.yml up -d --force-recreate

migrate:
	migrate create -ext sql -dir ./migrations -seq initial
