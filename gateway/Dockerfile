FROM golang:1.22-alpine as builder

# install protoc
RUN apk update && apk add protobuf

# install go dependencies
COPY ../go.mod ../go.sum /

# add proto files
COPY ../spec /spec

# add common
COPY ../common /common

# copy the source code
COPY ./gateway/src /app/src

# install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# generate the proto files
RUN protoc --proto_path=/spec --go_out=/spec --go-grpc_out=/spec --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative /spec/*.proto

WORKDIR /app

RUN go mod download

# build the binary
RUN go build -o /gateway src/*.go

# final stage
FROM gcr.io/distroless/base as final

COPY --from=builder /gateway /gateway

# server running port 
EXPOSE 7654

CMD ["/gateway"]
