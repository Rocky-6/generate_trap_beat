FROM golang:1.20-bullseye as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app .


FROM debian:bullseye-slim

COPY --from=builder /usr/local/bin/app /app
COPY --from=builder /usr/src/app/trap.db /trap.db

CMD [ "/app" ]