GOOS=linux GOARCH=386 go build
scp petra dev:petra
ssh dev sudo mv petra /bin/petra
