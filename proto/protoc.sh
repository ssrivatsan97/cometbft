cd tendermint/types
protoc -I=../ --gogo_out=. types.proto