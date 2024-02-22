build:
	docker build -t google-oauth2-user-service:latest .

tag:
	docker tag google-oauth2-user-service:latest docker.io/gamepkw/google-oauth2-user-service:latest

push:
	docker push gamepkw/google-oauth2-user-service:latest

stop:
	docker stop google-oauth2-user-service-container || true

remove:
	docker rm google-oauth2-user-service-container || true

run:
	docker run -d -p 9091:9091 --name google-oauth2-user-service-container google-oauth2-user-service:latest

#make build && make tag && make push && make stop && make remove && make run
