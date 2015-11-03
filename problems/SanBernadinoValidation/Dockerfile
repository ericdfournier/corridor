############################################################
# Dockerfile to build Go Language Container Image
# Based Upon Debian "Wheezy"
############################################################

FROM ericdfournier/corridor:latest
MAINTAINER Eric D Fournier: me@ericdfournier.com

WORKDIR /go/src/github.com/ericdfournier/corridor/problems/SanBernadinoValidation
RUN go get -u github.com/ericdfournier/corridor
RUN go run sanBernadinoValidation.go