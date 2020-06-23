FROM golang:1.14.3-alpine3.11 as builder
WORKDIR /go/baseweb/
COPY . .
RUN go build

FROM alpine:3.11
WORKDIR /go/
COPY --from=builder /go/baseweb/baseweb-simulation baseweb-simulation
COPY run.sh /root/run.sh
EXPOSE 8080
CMD /root/run.sh
