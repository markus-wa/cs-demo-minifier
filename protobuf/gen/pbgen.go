// Package gen contains the generated protobuf code.
// It is seperated from the protobuf package so it can be skipped for code coverage.
// Use 'go generate' to generate the code from the .proto files inside the proto sub directory.
package gen

// -I=proto is required, otherwise the generated .pb.go file will be put inside the proto directory.
// No idea what that is about to be honest . . .
//go:generate protoc -I=proto --gogofaster_out=. proto/replay.proto
