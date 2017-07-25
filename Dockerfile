FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/corest/giantswarmdemo
WORKDIR /go/src/github.com/corest/giantswarmdemo

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get
RUN go install 

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/giantswarmdemo
