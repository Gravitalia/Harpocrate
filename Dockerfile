FROM golang:1.20-alpine3.18 AS build

RUN apk update && apk add --no-cache gcc musl-dev

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux go build -o Harpocrate

FROM alpine:3.18 AS runtime

COPY --from=build /app/Harpocrate /app/Harpocrate

EXPOSE 5000
CMD [ "/app/Harpocrate" ]
