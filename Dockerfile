FROM golang:1.19 as builder

WORKDIR /app

RUN wget "https://github.com/go-swagger/go-swagger/releases/download/v0.30.3/swagger_linux_amd64" && \
    mv swagger_linux_amd64 /usr/local/bin/swagger && \
    chmod a+x /usr/local/bin/swagger

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY ./ /app

RUN wget "https://github.com/go-swagger/go-swagger/releases/download/v0.30.3/swagger_linux_amd64" && \
    mv swagger_linux_amd64 /usr/local/bin/swagger && \
    chmod a+x /usr/local/bin/swagger

RUN git config --global --add safe.directory /app

RUN make all

FROM alpine:3.17.0

RUN apk update && apk add tzdata

WORKDIR /app

COPY --from=builder /app/application-tracker-server /app/application-tracker-server

EXPOSE 8080

CMD [ "./application-tracker-server", "--port=8080", "--host=0.0.0.0" ]