

up: 
	docker-compose up

down: 
	docker-compose down

publisher:
	go run cmd/publisher/main.go

consumer:
	go run cmd/consumer/main.go


ip-inspect: 
	docker ps -q | xargs -n 1 docker inspect --format '{{ .Name }} {{range .NetworkSettings.Networks}} {{.IPAddress}}{{end}}' | sed 's#^/##' | sort -k 2


