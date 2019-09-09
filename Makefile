aws:
	GOOS=linux GOARCH=amd64 go build -o application
	zip -r aws.zip application migrations

docker:
	docker build -t "crawlyzer_auth" .

run:
	docker run -p 3000:3000 crawlyzer_auth

test:
	docker-compose up --exit-code-from test