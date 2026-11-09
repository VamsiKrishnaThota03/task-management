# ---- build ----
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- run ----
FROM alpine:3.22
RUN apk add --no-cache ca-certificates && adduser -D -H app
USER app
COPY --from=build /out/server /server

# Inside the container the server must listen on all interfaces; publish the
# port to 127.0.0.1 locally (docker run -p 127.0.0.1:8080:8080 ...).
ENV HOST=0.0.0.0 PORT=8080 GIN_MODE=release
EXPOSE 8080
CMD ["/server"]
