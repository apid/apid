# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
