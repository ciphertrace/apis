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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

func getToken(ctx context.Context, creds authData) (token oauth2.Token) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", creds.ClientID)
	form.Add("client_secret", creds.ClientSecret)
	form.Add("audience", creds.Audience[0])

	req, err := http.NewRequestWithContext(ctx, "POST", creds.TokenURI, bytes.NewBuffer([]byte(form.Encode())))
	if err != nil {
		panic(fmt.Sprintf("Could not create request for token: %v", err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Could not make request to get auth token: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Cannot read contents of response: %v", err))
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(fmt.Sprintf("Could not decode body because: %v\nBody was: %s", err, body))
	}

	return
}
