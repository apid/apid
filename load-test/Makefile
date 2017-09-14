docker:
	docker build -f Dockerfile.apid . -t apid-lt -t alexkhimich/apid-lt
dockerpush:
	docker push alexkhimich/apid-lt
dockerloadtest:
	docker rm apid-lt || true
	docker run -d --name apid-lt -e CLOUD_IP=localhost -p 9000:9000 apid-lt
clean:
	docker rmi -f apid-lt
	docker rmi -f alexkhimich/apid-lt
