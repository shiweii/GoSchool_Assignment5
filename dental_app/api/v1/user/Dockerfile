FROM golang:1.18

WORKDIR /code

RUN mkdir log

COPY pkg /pkg

COPY ../api/v1/user/go.mod ./
COPY ../api/v1/user/go.sum ./

RUN go mod download

COPY ../api/v1/user/*.go ./
COPY ../api/v1/user/.prod.env ./.env

RUN go build -o /main

CMD [ "/main" ]