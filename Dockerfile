FROM alpine:latest

# install necessaty packages
RUN apk add --no-cache vim git make musl-dev go

# set Env variables
# ?? maybe not needed ??

# add and change to user
RUN adduser -u 1001 -D user
WORKDIR /home/user
USER 1001

# install and build the blog server
RUN git clone https://github.com/Dolev123/goblog.git
WORKDIR /home/user/goblog
RUN go build


# setup the site
RUN mkdir /home/user/site 
RUN cp goblog /home/user/site
WORKDIR home/user/site
COPY ./config.json config.json
#COPY ./secrets.json ./.secrets.json

# set network
EXPOSE 8080/tcp

# start server
CMD ["goblog", "-config", "config.json"]

