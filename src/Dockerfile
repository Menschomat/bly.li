# Common base image with Go installation and necessary tools
FROM golang:1.23-alpine AS base-prepare

WORKDIR /src/

RUN adduser -S -u 10001 scratchuser
RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates

RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Argument for service name
ARG SERVICE

# Build stage for specific services
FROM base-prepare AS build

# Copy specific service files based on the argument
COPY services/${SERVICE}/ ./
COPY shared/ /shared/
RUN oapi-codegen -generate types,chi-server -package api -o api/api.gen.go api/openapi.yml

RUN CGO_ENABLED=0 go build -o /bin/blyli

# Final stage using scratch
FROM scratch

COPY --from=base-prepare /etc/passwd /etc/passwd
COPY --from=base-prepare /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base-prepare /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Europe/Berlin

# Copy the built binary from the build stage
COPY --from=build /bin/blyli /bin/blyli

USER scratchuser
ENTRYPOINT ["/bin/blyli"]
