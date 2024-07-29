FROM golang:1.22-alpine AS prepare

WORKDIR /src/

RUN adduser -S -u 10001 scratchuser
RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates

RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

FROM scratch AS base

COPY --from=prepare /etc/passwd /etc/passwd
COPY --from=prepare /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=prepare /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Europe/Berlin

##Blowup-Specific:
FROM prepare AS build

COPY services/blowup/ .
COPY shared/ /shared/
RUN oapi-codegen -generate types,chi-server -package api -o api/api.gen.go api/openapi.yml

RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM base

COPY --from=build /bin/blyli /bin/blyli

USER scratchuser
ENTRYPOINT ["/bin/blyli"]
