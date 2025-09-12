# Lista de servicios
SERVICES := achievement game metadata

# Puertos (ajusta si tus servicios usan otros)
PORT_achievement ?= 8082
PORT_game        ?= 8084
PORT_metadata    ?= 8083

# -------- Targets --------
.PHONY: build run run-all stop-all consul-up consul-down

# Compilar un servicio
build:
ifndef SERVICE
	$(error Usa: make build SERVICE=achievement|game|metadata)
endif
	go build -o bin/$(SERVICE) ./$(SERVICE)/cmd

# Compilar todos
build-all:
	@mkdir -p bin
	@for s in $(SERVICES); do \
		echo "==> Compilando $$s..."; \
		go build -o bin/$$s ./$$s/cmd || exit 1; \
	done

# Ejecutar un servicio
run:
ifndef SERVICE
	$(error Usa: make run SERVICE=achievement|game|metadata)
endif
	@PORT_VAR=PORT_$(SERVICE); \
	PORT=$${!PORT_VAR}; \
	echo "==> Ejecutando $(SERVICE) en puerto $$PORT..."; \
	go run ./$(SERVICE)/cmd -port $$PORT

# Ejecutar todos (en paralelo, cada uno en su puerto)
run-all:
	@echo "==> Ejecutando achievement en puerto $(PORT_achievement)"
	go run ./achievement/cmd -port $(PORT_achievement) &

	@echo "==> Ejecutando metadata en puerto $(PORT_metadata)"
	go run ./metadata/cmd -port $(PORT_metadata) &

	@echo "==> Ejecutando game en puerto $(PORT_game)"
	go run ./game/cmd -port $(PORT_game) &

	@echo "Servicios levantados en background:"
	@echo "  achievement -> http://localhost:$(PORT_achievement)/achievement"
	@echo "  metadata    -> http://localhost:$(PORT_metadata)/metadata"
	@echo "  game        -> http://localhost:$(PORT_game)/game"


# Detener todos (mata procesos 'go run')
stop-all:
	@pkill -f "go run" || true
	@echo "Servicios detenidos."

# Consul
consul-up:
	docker run -d --name consul-dev --rm -p 8500:8500 \
		hashicorp/consul:1.19 agent -dev -ui -client=0.0.0.0
	@echo "Consul levantado en http://localhost:8500/ui"

consul-down:
	-docker stop consul-dev
	@echo "Consul detenido"
