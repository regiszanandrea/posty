test:
	APP_ENV=testing ginkgo ./...
seed:
	APP_ENV=local go run database/seeder.go
