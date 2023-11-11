FROM golang:1.21-alpine AS build

WORKDIR /src/
COPY services/blowup/main.go services/blowup/go.* /src/
COPY shared/ /shared/
RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch
COPY --from=build /bin/blyli /bin/blyli
ENTRYPOINT ["/bin/blyli"]