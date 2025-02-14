.PHONY: help
help:
	@echo "build - build the project"
	@echo "clean - clean the project"
	@echo "run - run the project"
	@echo "dev - run the project with hot reload"
	@echo "tailwind-watch - watch the tailwind css"
	@echo "test - run tests"

.PHONY: build
build: templ tailwind-build
	go build -o bin/sws cmd/sws/main.go

.PHONY: clean
clean:
	rm -rf bin/sws web/static/css/output.css tmp/

.PHONY: dev
dev:
	air

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify

.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --watch

.PHONY: templ
templ:
	templ generate

.PHONY: test
test: test-unit test-integration

test-unit:
	go test -v ./...

test-integration:
	go test -v ./tests/integration

