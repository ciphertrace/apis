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
	"crypto/tls"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func getConn(token oauth2.Token, endpoint string) (connection *grpc.ClientConn) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS((&tls.Config{
			InsecureSkipVerify: true, // In production, you should add the necessary root CA's to your certificate store
		}))),
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(&token)),
	}

	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		panic(fmt.Sprintf("Could not create GRPC connection with endpoint %s because %v", endpoint, err))
	}

	return conn
}
