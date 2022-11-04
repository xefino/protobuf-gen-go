[![Unit Tests Status](https://github.com/xefino/protobuf-gen-go/actions/workflows/test.yml/badge.svg)](https://github.com/xefino/protobuf-gen-go/actions)

# protobuf-gen-go
This repository contains data types generated from files in Xefino's protobuf repository, that are intended to act as common dependencies to both public and private protobuf repositories.

## Guidelines

When developing code for this repository, the first few steps should be done inside the protobuf repository. However, once the code has been deployed to this repository, additional work may be required to fully develop the code resources necessary. This section contains guidelines pertinent to that extensibility.

### Generated Files

The files in this repository are generated. Therefore, any file with a `.pb.go` extension will be overwritten by subsequent releases. Therefore, under no circumstances should these files be changed. If changes are necessary, they can be done via the protobuf repository. Otherwise, extensions or utility functions may be written to add functionality as normal `.go` files will not be deleted.

### Releases

As this code is nearly entirely code-generated, updates to this repository are automatic. However, releases still need to be performed manually. Therefore, after changes have been pushed, please ensure that either a pre-release or production release is drafted so that the changes can be consumed by downstream services.

### Utils vs. Extensions

Packages in this repository typically contain one or both of a set of files, called `utils.go` and `extensions.go` respectively. The first of these, `utils.go`, is intended to include utility functions that would be useful for things like marshalling and unmarshalling. The other file, `extensions.go`, is intended to provide extra functionalilty that the standard protobuf files do not allow for. This might include things such as special representation code, conversion to other types, comparison, etc. The goal is for this functionality to be packaged together with the generated Go files so that data received from an RPC endpoint contains this functionality "out of the box".

### Serialization

The `utils/` directory contains a number of serialization helper functions that can be used to marshal enums to and from JSON, CSV, DynamoDB, SQL or other string-based formats. These are especially useful when data needs to be ingested or displayed in a separate format from the normal representation for a Go enum (int32).