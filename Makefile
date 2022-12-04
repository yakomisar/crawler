NAME = bin/crawler

$(NAME):
	@echo "Compiling for specific OS and Platform"
	go mod tidy
	go build -o bin/crawler main.go

run:
	go run main.go

clean:
	@echo "Deleting binary files..."
	@rm -rf bin/

all: $(NAME)