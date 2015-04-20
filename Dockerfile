############################################################
# Dockerfile to build Go Language Container Image
# Based Upon Debian "Wheezy"
############################################################

FROM golang:latest
MAINTAINER Eric D Fournier: me@ericdfournier.com
RUN go get github.com/ericdfournier/corridor