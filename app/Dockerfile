FROM okteto/okteto:stable AS okteto-cli

FROM golang:1.23-bookworm

COPY --from=okteto-cli /usr/local/bin/okteto /usr/local/bin/okteto

WORKDIR /app
COPY . .

RUN go build -o /usr/local/bin/app

CMD ["/usr/local/bin/app"]
