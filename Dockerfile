FROM golang:1.15.2 as builder

WORKDIR /go/src/app
ADD . /go/src/app

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app-golang

FROM gcr.io/distroless/static-debian10

USER nobody:nobody

COPY --from=builder /go/bin/app-golang /app-golang

EXPOSE 8080
CMD ["/app-golang"]
