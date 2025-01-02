BINPATH=bin/math-quiz

.PHONY: build
build: build-templ build-css build-js build-app

.PHONY: build-app
build-app:
	CGO_ENABLED=0 go build -o $(BINPATH) cmd/math-quiz/main.go

.PHONY: build-templ
build-templ:
	templ generate

.PHONY: build-css
build-css:
	npm --prefix web run build:css -- --minify

.PHONY: build-js
build-js:
	npm --prefix web run build:js -- --minify

#templ proxy localhost:7331
.PHONY: watch
watch:
	$(MAKE) -j4 watch-app watch-templ watch-css watch-js

.PHONY: watch-app
watch-app:
	go run github.com/air-verse/air@latest \
	--build.cmd "$(MAKE) build-app" \
	--build.bin "$(BINPATH)" \
	--build.include_ext "go" \
	--build.exclude_dir "bin, web"

.PHONY: watch-templ
watch-templ:
	templ generate \
	--watch \
	--proxy "http://localhost:3000" \
	--open-browser=false

.PHONY: watch-css
watch-css:
	npm --prefix web run build:css -- --watch=always

.PHONY: watch-js
watch-js:
	npm --prefix web run build:js -- --watch=forever

.PHONY: install-deps
install-deps:
	npm --prefix web install
	go mod download