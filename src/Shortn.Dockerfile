FROM golang:1.22-alpine AS build

WORKDIR /src/

RUN adduser -S -u 10001 scratchuser
RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates

COPY services/shortn/ .
COPY shared/ /shared/

RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Berlin

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/blyli /bin/blyli

USER scratchuser
ENTRYPOINT ["/bin/blyli"]