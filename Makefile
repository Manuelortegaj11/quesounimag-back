# Nombre del binario
BINARY_NAME=proyectoqueso

# Variables para facilitar los comandos
GO=go
GOFMT=gofmt
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOVET=$(GO) vet

.PHONY: all build test clean fmt vet run

# Primera regla all: La primera regla en el Makefile es all, y es la que se ejecutará por defecto cuando simplemente ejecutes make.
all: build

# La regla run depende de build
run: build
	@echo "Running the binary..."
	@./$(BINARY_NAME)

build: 
	@echo "Building the binary..."
	@$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)

drop-table:
	@echo "Dropping database"
	@cd cmd/database && go run main.go drop

migrate:
	@echo "Running migrations... "
	@cd cmd/database && go run main.go 

create-test-user:
	@echo "Creating a new user..."
	@cd cmd/database && go run main.go testuser
