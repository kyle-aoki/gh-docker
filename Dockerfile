FROM golang:1.22

EXPOSE 8080

COPY . .

RUN go build

CMD [ "./gh-docker" ]
