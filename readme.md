petra is for deploying docker tags that run stateless APIs or jobs.

APIs listen on port 8080 or 8081. nginx is used to reverse proxy requests from port 80 or 443 to ports 8080 or 8081.

```
build petra:

GOOS=linux GOARCH=386 go build
```

petra relies on "docker instances". A docker instance is an ec2 that runs a single container. 1 EC2 == 1 Container.

Docker instances are the simplest way to utilize containers. It avoids the devops overhead involved in managing ECS, EKS, or other container platforms.

A docker instance does not require petra to function. One can run the following commands to deploy new tags to docker instances:
```
ssh <instance>
docker pull <new-tag>
docker run -d -e CONFIG=$(cat config.json) -p <8080 or 8081>:8080 <new-tag>
<edit nginx config file to switch its reverse proxy from 8080 to 8081 (or vice versa)>
sudo nginx -s reload
docker stop <old-container>
docker rm <old-container>
```
Rather than run this sequence of commands every time, petra places these commands in a golang function. This function executes when petra detects a change in its config file (`/petra/petra-config.json`).

Programs that run via petra receive config through an environment variable called `CONFIG`.
