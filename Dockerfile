FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
ENTRYPOINT ["/bin/k9"]
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./LICENSE /LICENSE
COPY ./bin/k9-linux64 /bin/k9
