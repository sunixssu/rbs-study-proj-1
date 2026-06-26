drun:
	docker run -d -p 8080:8080 --name go-app go-service:latest

dbuild:
	docker build -t go-service-latest .