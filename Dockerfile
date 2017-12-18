FROM golang

# allow the user to pass in the GOPATH specific to their environment
ARG go_path=/go

# set the necessary environment variables for golang
ENV GOPATH $go_path
ENV PATH $GOPATH/bin:$PATH

# add the mailcave source to the Docker build directory
ADD . $GOPATH/src/github.com/tambchop/mailcave

# retrieve dependencies
WORKDIR $GOPATH/src/github.com/tambchop/mailcave/cmd/mailcave
RUN go get ./

# install the mailcave service
RUN go install github.com/tambchop/mailcave/cmd/mailcave

# the binary to run
ENTRYPOINT $GOPATH/bin/mailcave

# run this service on the following port
EXPOSE 8080
