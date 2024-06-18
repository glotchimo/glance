FROM node:18-alpine AS frontend-builder

WORKDIR /web

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build

FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=frontend-builder /web/dist ./web/dist

RUN go build -o /go-app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=backend-builder /go-app ./
COPY --from=backend-builder /app/web/dist ./web/dist

EXPOSE 8080

CMD ["./go-app"]
