drun:
	docker run -d \
	-v /proc/cpuinfo:/hostproc/cpuinfo:ro \
	-v /proc/meminfo:/hostproc/meminfo:ro \
	-p 8080:8080 \
	--name go-app \
	go-service:latest

dbuild:
	docker build -t go-service:latest .