FROM golang:1.21-alpine AS build

WORKDIR /src/
COPY services/shortn/main.go services/shortn/go.* /src/
COPY services/shortn/utils /src/utils/
COPY shared/ /shared/
RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch
COPY --from=build /bin/blyli /bin/blyli
ENTRYPOINT ["/bin/blyli"]