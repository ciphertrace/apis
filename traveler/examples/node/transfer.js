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
// use Brian's tools - var tools = require('./tools.js');
var debug = require('debug')('lookup');

// Root path from where the proto files are referenced
const ROOT_PROTO_PATH = process.env.ROOT_PROTO_PATH || pathLib.join(__dirname, "/../../..");

// Configuration - Set path to Traveler api.proto file
const TRAVELER_PROTO_PATH = process.env.TRAVELER_PROTO_PATH || pathLib.join(ROOT_PROTO_PATH, "/traveler/v1/api.proto");
// const AUTH_DATA = process.env.AUTH_DATA || '../credentials.json'; // Local credentials file with SECRET and CLIENT_ID
const AUTH_DATA = '../credentials.json';
const TRAVELER_ENDPOINT = process.env.TRAVELER_ENDPOINT || 'grpc.a639386.traveler.stage.cipheruse.com:443'; // CipherTrace Test

// Load the protocol buffers dynamically
const packageDefinition = protoLoader.loadSync(
  TRAVELER_PROTO_PATH,
  {
    keepCase: true, // was true
    longs: String,
    enums: String,
    defaults: true, // was true
//    bytes: Array,   // was not present
    arrays: true,
    oneofs: true,   // was true
    includeDirs: [ROOT_PROTO_PATH, pathLib.join(ROOT_PROTO_PATH, "/trisacrypto")]
  }
);

// Load the Traveler service package stub
const pb = grpc.loadPackageDefinition(packageDefinition).ciphertrace.apis.traveler.v1;

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
debug('-----------DEBUG-------------------');
debug('Create Credentials Object for gRPC endpoint');
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
debug('ROOT_PROTO_PATH',ROOT_PROTO_PATH);
debug('TRAVELER_ENDPOINT',TRAVELER_ENDPOINT);
debug('-------------------------------');
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



// make Transfer Request
function makeTransfer(client, transferdata) {
  client.transfer(transferdata, (err, response) => {

    debug('--------------------------------' );
    // debug('TransferResponse:', response);
    // debug(response);
    console.log(response);
    if (err) {
      console.log('Transfer Error:', err);
    }
  });
}

function getTransferData() {
var dateNow = new Date();

const msg2 = {
"information": {
  "identity": {
    "@type": "type.googleapis.com/ivms101.IdentityPayload",
    "originator": {
      "originator_persons": [
        {
          "natural_person": {
            "name": {
              "name_identifiers": [
                {
                  "primary_identifier": "Kelley",
                  "secondary_identifier": "George",
                  "name_identifier_type": "NATURAL_PERSON_NAME_TYPE_CODE_ALIA"
                }
              ]
            },
            "geographic_addresses": [
              {
                "street_name": "",
                "post_code": "",
                "country_sub_division":  "",
                "town_name": "",
                "country": ""
              }
            ],
            "national_identification": {
              "national_identifier": "",
              "national_identifier_type": "",
              "country_of_issue": "",
              "registration_authority": ""
            },
            "date_and_place_of_birth": {
              "date_of_birth": "",
              "place_of_birth": ""
            },
            "country_of_residence": ""
          }
        }
      ],
      "account_numbers": [
        ""
      ]
    },
    "beneficiary": {
        "beneficiary_persons": [
          {
            "natural_person": {
              "name": {
                "name_identifiers": [
                  {
                    "primary_identifier": "Jackson",
                    "secondary_identifier": "Peter",
                    "name_identifier_type": "NATURAL_PERSON_NAME_TYPE_CODE_ALIA"
                  }
                ]
              },
              "geographic_addresses": [
                {
                  "street_name": "",
                  "post_code": "",
                  "country_sub_division":  "",
                  "town_name": "",
                  "country": ""
                }
              ],
              "national_identification": {
                "national_identifier": "",
                "national_identifier_type": "",
                "country_of_issue": "",
                "registration_authority": ""
              },
              "date_and_place_of_birth": {
                "date_of_birth": "",
                "place_of_birth": ""
              },
              "country_of_residence": ""
            }
          }
        ],
        "account_numbers": [
          ""
        ]
      },
    "originating_vasp": {
      "originating_vasp": {
        "legal_person": {
          "name": {
            "name_identifiers": [
             {
                "legal_person_name": "Alice's Discount VASP, PLC",
                "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_LEGL"
              }
            ]
          },
          "geographic_addresses": [
            {
              "address_type": "",
              "street_name": "",
              "post_code": "",
              "town_name": "",
              "country": ""
            }
          ],
          "national_identification": {
            "national_identifier": "",
            "national_identifier_type": "",
            "country_of_issue": "",
            "registration_authority": ""
          },
          "country_of_registration": ""
        }
      }
    },
    "beneficiary_vasp": {
        "beneficiary_vasp": {
          "legal_person": {
            "name": {
              "name_identifiers": [
               {
                  "legal_person_name": "BobVASP",
                  "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_LEGL"
                }
              ]
            },
            "geographic_addresses": [
              {
                "address_type": "",
                "street_name": "",
                "post_code": "",
                "town_name": "",
                "country": ""
              }
            ],
            "national_identification": {
              "national_identifier": "",
              "national_identifier_type": "",
              "country_of_issue": "",
              "registration_authority": ""
            },
            "country_of_registration": ""
          }
        }
      },
    "transfer_path": {},
    "payload_metadata": {}
  },
  "transaction": {
    "@type": "type.googleapis.com/trisa.data.generic.v1beta1.Transaction",
    "txid": "",
    "originator":  "1ASkqdo1hvydosVRvRv2j6eNnWpWLHucMX",
    "beneficiary": "18nxAxBktHZDrMoJ3N2fk9imLX8xNnYbNh",
    "amount": "818.22",
    "network": "BTC",
    "assetType": "BTC",
    "timestamp": dateNow,
    "extra_json": "", 
    "tag": ""
  }
},  // end of information

crypto_address: {
            //crypto_address: "1LgtLYkpaXhHDu1Ngh7x9fcBs5KuThbSzw", // this is the Beneficiary VASP Address
            crypto_address: "1Hzej6a2VG7C8iCAD5DAdN72cZH5THSMt9", // this is the Beneficiary VASP Address
            network: "",
            asset_type: "",
            owner_name: "",
            owner_url: "",
            owner_country: "",
            owner_type: "",
            labels:[],
            vasp_id: "",               
            registered_directory:  "", 
            common_name: "",          
            endpoint: ""             
        },
contact_email: "",
envelope_id:   ""
};
debug('msg2 = ', msg2);
  return msg2;

}


function main() {
  const tokenPromise = getToken();
//string  jsonString = "";
//JsonFormat jsonFormat = new JsonFormat();
//jsonString = jsonFormat.printToString(IdentityPayload);

debug('-------------------------------');
debug('ROOT_PROTO_PATH',ROOT_PROTO_PATH);
debug('TRAVELER_PROTO_PATH',TRAVELER_PROTO_PATH);
debug('AUTH_DATA',AUTH_DATA);
debug('-------------------------------');
//debug(jsonString);
debug('-------------------------------');

  tokenPromise.then(response => {
    const accessToken = response.data.access_token;
    const client = getClient(accessToken);

//    checkStatus(client);

    const msg = getTransferData();
    makeTransfer(client, msg);

  });

}

main();