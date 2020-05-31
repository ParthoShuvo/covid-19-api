# build stage
FROM golang:1.13-alpine3.10 AS builder

# Add Maintainer info
LABEL maintainer="Shuvojit Saha"

# Set the current working directory inside the container 
WORKDIR /app

ADD . .
RUN go mod download 
RUN go build

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

ENV APP covid-19-api
ENV CONFIG covid_19_api.json
ENV DATA_FOLDER data
ENV APP_HOME /app

WORKDIR /root/
COPY --from=builder /app/$APP .
COPY --from=builder /app/$CONFIG ./$CONFIG
COPY --from=builder /app/data ./data

#PORT
EXPOSE 9000
# Start App
CMD ["./covid-19-api"]