
docker-build:
	docker build -t sickyoon/gohls:dev .

publish: docker-build
	docker push sickyoon/gohls:dev

