(1) example add env
export SSH_AUTH_LOGIN=login
export SSH_AUTH_KEY=`cat /Users/login/.ssh/id_rsa`

(2) example build and run
docker build . --tag=ssh-tunnel
docker run --env SSH_AUTH_LOGIN=$SSH_AUTH_LOGIN --env SSH_AUTH_KEY=$SSH_AUTH_KEY -p 8080:8080 -it ssh-tunnel --port=8080 --host=192.168.100.100 --server=server.domain

(3) example to hub
docker build . --tag username/ssh-tunnel:latest
docker push username/ssh-tunnel:latest

