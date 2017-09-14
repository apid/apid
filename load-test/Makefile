docker:
	docker build -f Dockerfile.apid . -t apid-lt -t alexkhimich/apid-lt
	docker build -f Dockerfile.mockserver . -t mockserver-lt
dockerpush:
	docker push alexkhimich/apid-lt
dockerloadtest:
	docker rm -f apid-lt || true
	docker rm -f mockserver-lt || true
	docker network rm apid-load-test || true
	docker network create apid-load-test
	docker run -d --name mockserver-lt -p 9001:9001 --network=apid-load-test mockserver-lt
	docker run -d --name apid-lt -p 9000:9000 --network=apid-load-test apid-lt
clean:
	docker rm -f apid-lt || true
	docker rm -f mockserver-lt || true
	docker rmi -f apid-lt || true
	docker rmi -f alexkhimich/apid-lt || true
	docker rmi -f mockserver-lt || true
	docker network rm apid-load-test || true
