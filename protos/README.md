## Install protoc
Download latest precompiled binary version of protoc from [protobuf repository](https://github.com/google/protobuf/releases). The file name starts with `protoc-`. 

Unzip the archive file, and simply place `bin/protoc` binary somewhere in your PATH.

Install protobuf plugin and grpc package.
```
$ go get -u google.golang.org/grpc
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

## Compile .proto files
If you are at the location of proto files, execute this command.
```
$ ./compile.sh
```