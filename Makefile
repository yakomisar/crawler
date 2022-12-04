crawler:
	echo "Program 'crawler'"

build:
	@echo "Compiling for specific OS and Platform"
	go build -o bin/crawler main.go

run:
	go run bin/crawler.go

clean:
	@echo "Deleting binary files..."
	@rm -rf bin/

all: crawler build