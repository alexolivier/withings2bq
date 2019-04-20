FROM golang:1.12-alpine AS builder
RUN apk add --no-cache git mercurial ca-certificates

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/alexolivier/withings2bq
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

FROM chrisjoyce911/goscratchcert
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app ./
ENTRYPOINT ["./app"]