FROM golang:1.19-alpine AS build

WORKDIR /src/
COPY blowup/main.go blowup/go.* /src/
COPY shared/ /shared/
RUN CGO_ENABLED=0 go build -o /bin/blyli

FROM scratch
COPY --from=build /bin/blyli /bin/blyli
ENTRYPOINT ["/bin/blyli"]