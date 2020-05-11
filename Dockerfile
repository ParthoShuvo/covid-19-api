FROM golang:1.13.10-alpine3.11

ENV APP covid-19-api
ENV APP_HOME /app
ENV DATASET_DOWNLOAD_SCRIPT covid-19-dataset.sh

RUN apk update
RUN apk upgrade
RUN apk add bash

# Remove previous APP_HOME Dir and Create Again
RUN rm -rf $APP_HOME; mkdir $APP_HOME

WORKDIR $APP_HOME/
ADD .  .

# REPLACE config.file -> uncomment following when localhost is binding with hostname
# OR change covid_19_api.json hostname (bind address) according to your specification
#======================================================================================
# ENV ORIGINAL_CONFIG_FILE covid_19_api.json
# ENV WITHOUT_HOST_CONFIG_FILE covid_19_api.no-host.json
# RUN cp -f $WITHOUT_HOST_CONFIG_FILE $ORIGINAL_CONFIG_FILE

# environment variable
RUN export GOFLAGS=-mod=vendor
RUN export GO111MODULE=on

# Fetching dependencies
RUN go mod download
RUN go mod vendor
RUN go mod verify

# Pulling Dataset
RUN chmod +x covid-19-dataset.sh
RUN ./covid-19-dataset.sh

#Build The binary
RUN go build

#PORT
EXPOSE 9000

# Start App
CMD [ "./covid-19-api" ]