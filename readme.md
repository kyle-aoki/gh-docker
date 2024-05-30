petra is for deploying docker tags that run stateless APIs or jobs.

APIs listen on port 8080 or 8081. nginx is used to reverse proxy requests from port 80 or 443 to ports 8080 or 8081.

```
build petra:

GOOS=linux GOARCH=386 go build
```
