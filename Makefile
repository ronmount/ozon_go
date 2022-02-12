PG = postgresql
REDIS = redis

all:
		@echo "Specify storage type: make postgresql or make redis"

test:
		@echo "Running tests...."
		@go test ./internal/config_parser
		@go test ./internal/tools

$(PG):
		docker-compose --profile postgresql up -- build


$(REDIS):
		docker-compose --profile redis up --build