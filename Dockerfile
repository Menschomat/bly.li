FROM golang:1.19-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
COPY utils /src/utils/
COPY model /src/model/
COPY redis /src/redis/
RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch
COPY --from=build /bin/blyli /bin/blyli
ENTRYPOINT ["/bin/blyli"]