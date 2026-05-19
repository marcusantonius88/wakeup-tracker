FROM golang:1.22-alpine AS build

WORKDIR /src
COPY backend/go.mod ./
COPY backend/shared ./shared
COPY backend/services ./services

ARG SERVICE
RUN go build -o /out/service ./services/${SERVICE}/cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/service /app/service
EXPOSE 8080
ENTRYPOINT ["/app/service"]

