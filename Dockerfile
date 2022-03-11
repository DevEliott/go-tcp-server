FROM golang:latest AS build

WORKDIR /app

COPY . .

# CGO_ENABLED=0 so its statically link exec and -ldflags="-w -s" to remove debug infos
RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-w -s" -a -o ./main .

FROM scratch

COPY --from=build ./app/main /

EXPOSE 8080

ENTRYPOINT [ "/main" ]