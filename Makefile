build:
	export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@3/lib/pkgconfig" && go build -tags dynamic .
run:
	source .env && SERVICE_ID=http && go run main.go