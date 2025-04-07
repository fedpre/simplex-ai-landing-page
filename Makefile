.PHONY: dev
dev:
	air

.PHONY: build
build:
	go build -o main .

.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch --minify
