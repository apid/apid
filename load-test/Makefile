docker:
	docker build -f Dockerfile.apid . -t apid-lt -t alexkhimich/apid-lt
	docker build -f Dockerfile.mockedcloud . -t apid-lt-mockedcloud -t alexkhimich/apid-lt-mockedcloud
dockerpush:
	docker push alexkhimich/apid-lt
	docker push alexkhimich/apid-lt-mockedcloud
dockerloadtest:
	docker run -d --name apid-lt-mockedcloud -e CLOUD_IP=192.168.99.100 -p 9001:9001 apid-lt-mockedcloud
	docker run --name apid-lt -e CLOUD_IP=192.168.99.100 -p 9000:9000 apid-lt
