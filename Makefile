docker:
	docker build -f load-test/Dockerfile.apid ./load-test -t apid-lt -t alexkhimich/apid-lt
	docker build -f load-test/Dockerfile.mockedcloud ./load-test -t apid-lt-mockedcloud -t alexkhimich/apid-lt-mockedcloud
dockerpush:
	docker push alexkhimich/apid-lt
	docker push alexkhimich/apid-lt-mockedcloud
