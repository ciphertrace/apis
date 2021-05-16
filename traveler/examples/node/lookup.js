// Copyright 2021 CipherTrace Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const axios = require('axios').default;
const fs = require('fs');
const qs = require('qs');
const pathLib = require("path");

// Root path from where the proto files are referenced
const ROOT_PROTO_PATH = process.env.ROOT_PROTO_PATH || pathLib.join(__dirname, "/../../../..");

// Configuration - Set path to Traveler api.proto file
const TRAVELER_PROTO_PATH = process.env.TRAVELER_PROTO_PATH || pathLib.join(ROOT_PROTO_PATH, "/apis/traveler/v1/api.proto");
const AUTH_DATA = process.env.AUTH_DATA || '../credentials.json'; // Local credentials file with SECRET and CLIENT_ID
const TRAVELER_ENDPOINT = process.env.TRAVELER_ENDPOINT || 'grpc.a639386.traveler.stage.cipheruse.com:443'; // CipherTrace Test

// Load the protocol buffers dynamically
const packageDefinition = protoLoader.loadSync(
  TRAVELER_PROTO_PATH,
  {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
    includeDirs: [ROOT_PROTO_PATH]
  }
);

// Load the Traveler service package stub
const pb = grpc.loadPackageDefinition(packageDefinition).ciphertrace.singularity.apis.traveler.v1;

// Load the authentication json file
const rawData = fs.readFileSync(AUTH_DATA);
const authData = JSON.parse(rawData);

// Authentication with OAuth2 - get access token
function getToken() {
  const tokenpromise = axios.request({
    method: 'POST',
    url: authData.token_uri,
    headers: {'content-type': 'application/x-www-form-urlencoded'},
    data: qs.stringify({
      'grant_type': 'client_credentials',
      'client_id': authData.client_id,
      'client_secret': authData.client_secret,
      'audience': authData.audience[0]
    })
  });
  return tokenpromise;
}

// Get client to make gRPC calls
function getClient(accessToken) {
  try {
    // Create the credentials object to provide to gRPC endpoint
    const channelCreds = grpc.credentials.createSsl();
    const metaCallback = (_params, callback) => {
      const meta = new grpc.Metadata();
      // meta.add('x-endpoint-api-userinfo', 'Bearer ' + accessToken);
      meta.add('Authorization', 'Bearer ' + accessToken);
      callback(null, meta);
    }
    const callCreds = grpc.credentials.createFromMetadataGenerator(metaCallback);
    const combCreds = grpc.credentials.combineChannelCredentials(channelCreds, callCreds);

    // Create the gRPC client to run gRPC calls pointing to the given endpoint with credentials containing access token 
    const client = new pb.Traveler(TRAVELER_ENDPOINT, combCreds);
    return client;
  } catch (e) {
    console.log(e);
    throw e;
  }
}

// Check status of gRPC endpoint
function checkStatus(client) {
  client.status({no_stream_info: true}, (err, response) => {
    if (err) {
      console.log('Error:', err)
    }
    console.log('Status:', response);
  });
}

// Check crypto addresses
function checkAddresses(client, addresses) {
  client.lookupAddress(addresses, (err, response) => {
    console.log('LookupResponse:', response);
    if (err) {
      console.log('Lookup Error:', err);
    }
  });
}

function getLookupAddresses() {
  const addresses = {
    'addresses': [{'crypto_address': 'LQcctH7yCV1qh97ktAHWHAbDHtc5hgCdAF', 'network': 'ltc_wallet'},
      {'crypto_address': 'pzt7ezr9gpdph6n283j657fc4qscrdwx3c9hdpaldy', 'network': 'bch_wallet'},
      {'crypto_address': '17aaJMZqArjH3eufDaEJ3HZg74BJT7Gmgz', 'network': 'bch_wallet'},
      {'crypto_address': 'bnb12wpk84anfltnwj9l7l9xazur4u70yxmdyj8glk', 'network': 'bnb_wallet'}]
  };
  return addresses;
}

function main() {
  const tokenPromise = getToken();

  tokenPromise.then(response => {
    const accessToken = response.data.access_token;
    const client = getClient(accessToken);

    checkStatus(client);

    const addresses = getLookupAddresses();
    checkAddresses(client, addresses);
  });
}

main();
