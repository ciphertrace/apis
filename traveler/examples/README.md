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
* ROOT_PROTO_PATH: Path to the root of the protocol buffer specifications needed by Traveler.
  
* TRAVELER_PROTO_PATH: Path to the Traveler protocol buffer specification.
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
