# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.6

ENV GO15VENDOREXPERIMENT="1"

LABEL sponge.version="0.1"
LABEL sponge.usage="docker run -d sponge  --env-file=PATHTOENVFILE /bin/sh -c \"sponge import\""
LABEL sponge.license="MIT"

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/coralproject/sponge

RUN echo $GOPATH

# Build & Install
RUN cd /go/src/github.com/coralproject/sponge/cmd/sponge && go install
#go build github.com/coralproject/sponge/cmd/sponge/main.go

ENV PATH /go/bin:$PATH

# Run the app
CMD ["sponge", "import"]
