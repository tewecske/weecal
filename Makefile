run: build
	@./bin/app

.PHONY: build
build:
	@test -f web/public/htmx.min.js || wget -q -O web/public/htmx.min.js https://unpkg.com/htmx.org@2.0.1
	@echo -n 'web/public/htmx.min.js sha256-'; cat web/public/htmx.min.js | openssl sha256 -binary | openssl base64
	@test -f web/public/response-targets.js || wget -q -O web/public/response-targets.js https://unpkg.com/htmx-ext-response-targets@2.0.0/response-targets.js
	@echo -n 'web/public/response-targets.js sha256-'; cat web/public/response-targets.js | openssl sha256 -binary | openssl base64
	@go build -tags dev -o bin/app cmd/app/main.go

.PHONY: templ
templ:
	templ generate --watch --proxy=http://localhost:3000

.PHONY: css
css:
	tailwindcss -i web/css/app.css -o web/public/styles.css --watch

.PHONY: clean
clean:
	rm -rf tmp/
	rm web/public/styles.css
	find web/templates/ -name '*templ.go' -type f -print0 | xargs -0 rm

.PHONY: test
test:
	  go test -race -v -timeout 30s ./...
