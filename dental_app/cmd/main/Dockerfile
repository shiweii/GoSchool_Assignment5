FROM golang:1.18

WORKDIR /code

RUN mkdir log
RUN mkdir templates

COPY pkg /pkg

COPY ../cmd/main/go.mod ./
COPY ../cmd/main/go.sum ./

RUN go mod download

COPY ../cmd/main/*.go ./
COPY ../cmd/main/.prod.env ./.env
COPY ../cmd/main/templates/*.* ./templates


RUN go build -o /main

CMD [ "/main" ]