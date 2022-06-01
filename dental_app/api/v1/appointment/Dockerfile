FROM golang:1.18

WORKDIR /code

RUN mkdir log

COPY pkg /pkg

COPY ../api/v1/appointment/go.mod ./
COPY ../api/v1/appointment/go.sum ./

RUN go mod download

COPY ../api/v1/appointment/*.go ./
COPY ../api/v1/user/.prod.env ./.env

RUN go build -o /main

CMD [ "/main" ]