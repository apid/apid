docker:
	docker build -f Dockerfile.apid . -t apid-lt -t alexkhimich/apid-lt
	docker build -f Dockerfile.mockserver . -t mockserver-lt
dockerpush:
	docker push alexkhimich/apid-lt
dockerloadtest:
	docker rm apid-lt || true
	docker rm mockserver-lt || true
	docker run -d --name mockserver-lt -p 9001:9001 mockserver-lt
	docker run -d --name apid-lt -e CLOUD_IP=localhost -p 9000:9000 apid-lt
clean:
	docker rmi -f apid-lt || true
	docker rmi -f alexkhimich/apid-lt || true
	docker rmi -f mockserver-lt || true
