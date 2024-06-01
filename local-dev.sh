GOOS=linux GOARCH=386 go build
scp petra dev.student-api.kyle-aoki.dev:petra
ssh dev.student-api.kyle-aoki.dev sudo mv petra /bin/petra
