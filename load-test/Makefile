docker:
	docker build -f Dockerfile.apid . -t apid-lt -t alexkhimich/apid-lt
	docker build -f Dockerfile.mockserver . -t mockserver-lt
dockerpush:
	docker push alexkhimich/apid-lt
dockerloadtest:
	docker rm apid-lt || true
	docker rm mockserver-lt || true
	docker run -d --name mockserver-lt -p 9001:9001 --network=bridge mockserver-lt
	docker run -d --name apid-lt -e CLOUD_IP=$(docker network inspect bridge | grep "Gateway" | cut -d \" -f 4) -p 9000:9000 --network=bridge apid-lt
clean:
	docker rm -f apid-lt || true
	docker rm -f mockserver-lt || true
	docker rmi -f apid-lt || true
	docker rmi -f alexkhimich/apid-lt || true
	docker rmi -f mockserver-lt || true
