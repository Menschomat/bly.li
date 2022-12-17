FROM golang:1.19-alpine AS build

WORKDIR /src/
COPY shortn/main.go shortn/go.* /src/
COPY shortn/utils /src/utils/
COPY shared/ /shared/
RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch
COPY --from=build /bin/blyli /bin/blyli
ENTRYPOINT ["/bin/blyli"]