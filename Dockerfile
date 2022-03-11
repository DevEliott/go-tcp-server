FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN go mod tidy && CGO_ENABLED=0 go build -a -o ./main .

FROM scratch

COPY --from=build ./app/main /

EXPOSE 8080

ENTRYPOINT [ "/main" ]