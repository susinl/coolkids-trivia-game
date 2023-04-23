ARG BUILDER_IMAGE=golang:buster
ARG BASE_IMAGE=alpine:3.15.4

################Builder Image###################
FROM $BUILDER_IMAGE AS builder

WORKDIR /app
COPY . .

RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp ./cmd/main.go

###############Application Image################
FROM $BASE_IMAGE as final

ENV TZ=Asia/Bangkok
WORKDIR /app
COPY --from=builder /app/goapp .

CMD ["/app/goapp"]