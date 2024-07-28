FROM golang:1.22-alpine AS build

WORKDIR /src/

RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates

COPY services/shortn/main.go services/shortn/go.* /src/
COPY services/shortn/utils /src/utils/
COPY shared/ /shared/

RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/blyli /bin/blyli
ENV TZ Europe/Berlin
ENTRYPOINT ["/bin/blyli"]