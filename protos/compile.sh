#!/bin/bash
PRJ_ROOT_DIR="$(cd "$(dirname "$0")"/.. && pwd)"
PROTO_ROOT_DIR=${PRJ_ROOT_DIR}/protos

for protos in $(find "$PROTO_ROOT_DIR" -name '*.proto' -exec dirname {} \; | sort | uniq) ; do
    protoc --proto_path="$PROTO_ROOT_DIR" --go_out=plugins=grpc,paths=source_relative:"$PRJ_ROOT_DIR"/protos "$protos"/*.proto
done