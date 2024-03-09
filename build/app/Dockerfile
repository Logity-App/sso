FROM golang:1.22-alpine as builder

WORKDIR /build
COPY go.mod .
RUN go mod download
RUN go mod tidy
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /sso cmd/sso/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder sso /bin/sso

ENTRYPOINT ["/bin/sso"]