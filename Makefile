generate:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/admin/admin.proto

api:
	docker build -t ghcr.io/0fau/logs-api:latest -f pkg/cmd/api/Dockerfile .
	docker push ghcr.io/0fau/logs-api:latest

web:
	docker build -t ghcr.io/0fau/logs:latest -f web/Dockerfile web
	docker push ghcr.io/0fau/logs:latest

bot:
	docker build -t ghcr.io/0fau/logs-bot:latest -f pkg/cmd/bot/Dockerfile .
	docker push ghcr.io/0fau/logs-bot:latest

admin:
	docker build -t ghcr.io/0fau/logs-admin:latest -f pkg/cmd/admin/Dockerfile .
	docker push ghcr.io/0fau/logs-admin:latest

screenshot:
	docker build -t ghcr.io/0fau/logs-screenshot:latest -f pkg/cmd/screenshot/Dockerfile .
	docker push ghcr.io/0fau/logs-screenshot:latest