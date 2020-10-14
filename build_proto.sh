#!/usr/bin/env bash

protoc -I onebot_idl --gofast_out=proto_gen/onebot onebot_idl/*.proto
protoc -I dto_proto --gofast_out=proto_gen/dto dto_proto/*.proto