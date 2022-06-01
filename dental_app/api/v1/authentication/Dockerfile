FROM golang:1.18

WORKDIR /code

RUN mkdir log

COPY pkg /pkg

COPY ../api/v1/authentication/go.mod ./
COPY ../api/v1/authentication/go.sum ./

RUN go mod download

COPY ../api/v1/authentication/*.go ./
COPY ../api/v1/authentication/.prod.env ./.env

RUN go build -o /main

CMD [ "/main" ]