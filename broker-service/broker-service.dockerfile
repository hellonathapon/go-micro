# # base go image
# FROM golang:1.19-alpine as builder

# # create /app dir
# RUN mkdir /app

# # copy this directory into docker file system /app
# COPY . /app

# # set working directory to /app
# WORKDIR /app

# # build the go source code there
# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# # add executable priviledge to brokerApp directory
# RUN chmod +x /app/brokerApp

# ----------------------------------------------------------------------------------------------------------------------------#
#   I use `Make` to do compile Go source code on local machine and therefore there is no need to                              #
#   copy the Go source code and compile it again in docker image, Just copy the executable binary file to the docker image.   #
# ----------------------------------------------------------------------------------------------------------------------------#

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

# copy built go.exe from brokerApp of base image to this tiny image /app
COPY --from=builder /app/brokerApp /app

# run the executable compiled file in this image without the go runtime/compiler image
# this makes the image more efficient and faster to run
CMD ["/app/brokerApp"]