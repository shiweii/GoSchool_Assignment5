module dental_app/cmd/main

go 1.18

require (
	github.com/gorilla/mux v1.8.0
	github.com/shiweii/logger v0.0.0-00010101000000-000000000000
	github.com/shiweii/utility v0.0.0-00010101000000-000000000000
	github.com/shiweii/validator v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace (
	github.com/shiweii/logger => ../../pkg/logger
	github.com/shiweii/utility => ../../pkg/utility
	github.com/shiweii/validator => ../../pkg/validator
)
