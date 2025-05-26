run:
	go run ./cmd/app
	
docker_run_dev: 
	docker compose -f deploy/dev/docker-compose.yml -p ai-chat-dev up -d
docker_stop_dev:
	docker compose -f deploy/dev/docker-compose.yml -p ai-chat-dev down

migrate_create:
	migrate create -ext sql -dir infrastructure/driven/repository/db/migration -seq init_schema
migrate_up:
	migrate -path infrastructure/driven/repository/db/migration -database "postgresql://root:password@localhost:5432/ai_chatbot?sslmode=disable" -verbose up
migrate_down: 
	migrate -path infrastructure/driven/repository/db/migration -database "postgresql://root:password@localhost:5432/ai_chatbot?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: run docker_run_dev docker_stop_dev migrate_create migrate_up migrate_down sqlc