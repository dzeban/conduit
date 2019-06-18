FROM golang:1.12-alpine as builder
WORKDIR /conduit/
COPY . .
RUN apk --no-cache add make
RUN make install

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/conduit /conduit
EXPOSE 8080
CMD ["/conduit"]
