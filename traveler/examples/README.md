# CipherTrace Traveler API Examples

## Getting Started
### Credentials
These examples will need the staging credentials provided to your
organization by your CipherTrace representative.
They assume these credentials are stored in a
_credentials.json_ file in the root of the examples directory.

### Configuration
Secondly, these examples can be customized using the following environment
variables:
* ROOT_PROTO_PATH: Path to the root of the protocol buffer specifications needed by Traveler. (Node only)
* TRAVELER_PROTO_PATH: Path to the Traveler protocol buffer specification. (Node only)
* AUTH_DATA: Path to the _credentials.json_ file
* **TRAVELER_ENDPOINT**: Traveler endpoint provided to you by your CipherTrace representative.

> Note:
> The path variables will default to the correct location when the examples are run from their default directory.

An example file with empty and commented out variables is provided for your convenience.
Copy the `sample.env` file to `.env`, remove the `#` symbol in front of the variable
that you want to specify and provide its value after the `=` sign.

# Running the examples
## Node
The node examples are provided under the `node` directory.

In a terminal window go to the `examples/node` directory and run:
```shell
npm install
```
To install all dependencies.

After that scripts can be run using node directly or using the `npm run` command
if you want the configuration to be loaded from the `.env` file in the root:
```shell
npm run lookup
```

## Go
### Prerequisites
In order to run the Go example you will need the Go compiler, the protobuf compiler, the go command for generating protobuf outputs and the go command for creating grpc outputs. Please follow the directions in the "Prerequisites" section on this page:
https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers

### Building
In a terminal window go to the `examples/go` directory and run:
```shell
go generate ./...
go build -o client .
```

### Running
To run the client after building, be sure you are still in the `examples/go` directory and then you can run this to test the connection and status of your traveler node:
```shell
./client
```
To test a lookup (reading from the `lookup.json` file in the examples directory) just run:
```shell
./client -lookup
```
And finally to test a transfer (using the data from the `transfer.json` file in the examples directory) run:
```shell
./client -transfer
```