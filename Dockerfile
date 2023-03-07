FROM golang:1.19.6-alpine AS builder

ENV APP_HOME /go/src/jobinterviewerapp

WORKDIR "$APP_HOME"
COPY ./ ./

RUN go mod download
RUN go mod verify
RUN go build -o jobinterviewerapp  ./cmd/telegram/main.go

FROM golang:1.19.6-alpine

ENV APP_HOME /go/src/jobinterviewerapp
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME"/jobinterviewerapp $APP_HOME

EXPOSE 8010
CMD ["./jobinterviewerapp"]
