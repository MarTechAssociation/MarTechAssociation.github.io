build:
	export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@3/lib/pkgconfig" && go build -tags dynamic .
run:
	source .env && source .env_dev && SERVICE_ID=http && go run main.go
gen:
	curl -XGET "http://localhost:8080/gen"
change-ruby:
	chruby 3.2.2