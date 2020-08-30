#!/bin/bash
# docker-compose does not work on armhf
# so this script recreates the docker container
docker ps|grep "logvac_dev"|cut -d " " -f1|xargs docker rm -f
docker run --name logvac_dev -p 1514:1514/udp -p 6360-6361:6360-6361 -d thordin9/logvac
echo 'Now GET /add-token with headers X-USER-TOKEN and X-AUTH-TOKEN e.g. curl -H "X-AUTH-TOKEN: secret" -H "X-USER-TOKEN: secret" http://127.0.0.1:6360/add-token -v'
echo 'then /logs?type=app with headers X-USER-TOKEN and Content-Type for application/json e.g. curl -H "Content-Type: application/json" -H "X-USER-TOKEN: secret" http://127.0.0.1:6360/logs?type=app -v'
