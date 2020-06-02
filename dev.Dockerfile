# build stage
FROM golang:1.13-alpine3.10 AS builder

# Add Maintainer info
LABEL maintainer="Shuvojit Saha"

# Set the current working directory inside the container 
WORKDIR /app

ADD . .

# REPLACE .env with .env.Docker
ENV ENV .env
ENV ENV_DOCKER .env.Docker
RUN cp -f $ENV_DOCKER $ENV

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
ENV ENV_FILE .env

WORKDIR /root/
COPY --from=builder /app/$APP .
COPY --from=builder /app/$CONFIG ./$CONFIG
COPY --from=builder /app/data ./data
COPY --from=builder /app/$DATASET_SCRIPT ./$DATASET_SCRIPT
COPY --from=builder /app/$ENV_FILE ./$ENV_FILE

# Make dataset downloader script executable
RUN chmod +x $DATASET_SCRIPT

#PORT
EXPOSE 9000
# Start App
CMD ["sh", "-c", "./${DATASET_SCRIPT} && ./covid-19-api"]