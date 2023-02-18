FROM alpine:latest

RUN mkdir /app

# copy built go.exe from brokerApp of base image to this tiny image /app
COPY loggerServiceApp /app

# run the executable compiled file in this image without the go runtime/compiler image
# this makes the image more efficient and faster to run
CMD ["/app/loggerServiceApp"]