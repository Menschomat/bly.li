FROM golang:1.22-alpine AS build

WORKDIR /src/

RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates

RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

COPY services/blowup/ .
COPY shared/ /shared/
RUN oapi-codegen -generate types,chi-server -package api -o api/blowup.gen.go api/openapi.yml

# Build the Go app
RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Berlin

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/blyli /bin/blyli

ENTRYPOINT ["/bin/blyli"]