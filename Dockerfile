FROM okteto/okteto:stable AS okteto-cli

FROM golang:1.22.7-bookworm

COPY --from=okteto-cli /usr/local/bin/okteto /usr/local/bin/okteto

WORKDIR /app
ADD deployer/ .
RUN go build -o /usr/local/bin/deployer

CMD ["/usr/local/bin/deployer"]