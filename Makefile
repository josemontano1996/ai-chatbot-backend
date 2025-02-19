run:
	go run ./cmd/app
docker_run_dev: 
	docker compose -f deploy/dev/docker-compose.yml -p ai-chat-dev up -d
docker_stop_dev:
	docker compose -f deploy/dev/docker-compose.yml -p ai-chat-dev down

.PHONY: run docker_run_dev docker_stop_dev