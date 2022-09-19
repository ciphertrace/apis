package main

// Copyright 2022 CipherTrace Inc
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

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	common "github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1"
	api "github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/v1"
	data "github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/data/generic/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	ctx, travelerEndpoint, credentialFile, doLookup, doTransfer := parseOptions()

	creds := readAuthData(credentialFile)

	token := getToken(ctx, creds)

	conn := getConn(token, travelerEndpoint)
	defer conn.Close()

	client := api.NewTravelerClient(conn)

	status, err := client.Status(ctx, &api.StatusRequest{NoStreamInfo: true, NoConciergeInfo: true, NoSunriseInfo: true})
	if err != nil {
		log.Fatalf("Failed to get status from server %s because %v", travelerEndpoint, err)
		return
	}

	if status.Health != "ok" {
		log.Fatalf("Health of the node is not ok, full status is %++v", status)
		return
	}

	var outPretty protojson.MarshalOptions
	outPretty.Multiline = true
	outPretty.Indent = "\t"
	outPretty.EmitUnpopulated = false

	prettyStatus, _ := outPretty.Marshal(status)
	log.Printf("Successfully got status from server, status is: %s", prettyStatus)

	if doLookup {
		performLookup(ctx, client)
	}

	if doTransfer {
		performTransfer(ctx, client)
	}

	fmt.Println("Client executed successfully")
}

func performLookup(ctx context.Context, client api.TravelerClient) {
	f, err := os.Open("../lookup.json")
	if err != nil {
		panic(fmt.Sprintf("Cannot open ../lookup.json: %v", err))
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	var jsonAddresses struct {
		Addresses []*common.CryptoAddress `json:"addresses"`
	}

	err = dec.Decode(&jsonAddresses)
	if err != nil {
		panic(fmt.Sprintf("Could not parse ../lookup.json: %v", err))
	}

	reply, err := client.LookupAddress(ctx, &api.LookupAddressRequest{Addresses: jsonAddresses.Addresses})
	if err != nil {
		log.Printf("Failed to lookup address, got error: %v", err)
		return
	}

	prettyOut, _ := json.MarshalIndent(reply, "", "\t")
	fmt.Printf("Lookup returned the following results: \n%s\n", prettyOut)
}

func performTransfer(ctx context.Context, client api.TravelerClient) {
	f, err := os.Open("../transfer.json")
	if err != nil {
		panic(fmt.Sprintf("Cannot open ../transfer.json: %v", err))
	}
	jsonBytes, _ := io.ReadAll(f)
	f.Close()

	var xfer api.TransferRequest

	var parser protojson.UnmarshalOptions
	err = parser.Unmarshal(jsonBytes, &xfer)
	if err != nil {
		panic(fmt.Sprintf("Could not parse ../transfer.json: %v", err))
	}

	var txn data.Transaction
	txnRaw := xfer.Request.(*api.TransferRequest_Information).Information.Transaction

	err = txnRaw.UnmarshalTo(&txn)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal the transaction in the information message: %v", err))
	}
	txn.Timestamp = time.Now().Format(time.RFC3339)
	txnRaw.MarshalFrom(&txn)

	reply, err := client.Transfer(ctx, &xfer)
	if err != nil {
		log.Printf("Failed to perform a transfer: %v", err)
		log.Printf("Transfer request was: %++v", &xfer)
		return
	}

	var outPretty protojson.MarshalOptions
	outPretty.Multiline = true
	outPretty.Indent = "\t"
	outPretty.EmitUnpopulated = false

	prettyOut, _ := outPretty.Marshal(reply)
	fmt.Printf("Transfer successful, returned the following response: \n%s\n", prettyOut)
}

func parseOptions() (ctx context.Context, travelerEndpoint, credentialFile string, doLookup, doTransfer bool) {
	travelerEndpoint = os.Getenv("TRAVELER_ENDPOINT")
	if travelerEndpoint == "" {
		travelerEndpoint = "grpc.a639386.traveler.stage.cipheruse.com:443"
	}

	credentialFile = os.Getenv("AUTH_DATA")
	if credentialFile == "" {
		credentialFile = "../credentials.json"
	}

	flag.BoolVar(&doLookup, "lookup", false, "Perform a lookup using the example data from ../lookup.json")
	flag.BoolVar(&doTransfer, "transfer", false, "Initiate a transfer request using the example data from ../transfer.json")

	flag.Parse()

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
			signal.Stop(sigChan)
		}
	}()

	return
}

type authData struct {
	TokenURI     string   `json:"token_uri"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Audience     []string `json:"audience"`
}

func readAuthData(filePath string) (data authData) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Cannot open file %s: %v", filePath, err))
	}
	defer f.Close()

	j := json.NewDecoder(f)
	err = j.Decode(&data)

	if err != nil {
		panic(fmt.Sprintf("Cannot JSON parse file %s: %v", filePath, err))
	}

	return
}
