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
# Add /bin/bash
RUN apk add --no-cache bash

ENV APP covid-19-api
ENV CONFIG covid_19_api.json
ENV DATA_FOLDER data
ENV APP_HOME /app
ENV DATASET_SCRIPT covid-19-dataset.sh

WORKDIR /root/
COPY --from=builder /app/$APP .
COPY --from=builder /app/$CONFIG ./$CONFIG
COPY --from=builder /app/data ./data
COPY --from=builder /app/$DATASET_SCRIPT ./$DATASET_SCRIPT

# Make dataset downloader script executable
RUN chmod +x $DATASET_SCRIPT

#PORT
EXPOSE 9000
# Start App
CMD ["sh", "-c", "./${DATASET_SCRIPT} && ./covid-19-api"]