package main

// Please note: This is not the recommended way to compile the protobufs.
// There are much better options to do this than manually building the compile scripts by hand, look at bazel for example as a good option to auto generate these protobufs in a more managable and maintainable way.
// This example code was built this way to enable people to try out this code without having to install a separate build system to run it.

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

//go:generate -command buildproto protoc -I ../../../trisacrypto/ -I ../../../ --go_out=pb/ --proto_path=../../../ --go_opt=module=github.com/ciphertrace/apis/traveler/examples/go/pb --go_opt=Mivms101/enum.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/ivms101 --go_opt=Mivms101/ivms101.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/ivms101 --go_opt=Mivms101/identity.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/ivms101 --go_opt=Mtrisa/api/v1beta1/errors.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/api/v1beta1 --go_opt=Mtrisa/api/v1beta1/api.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/api/v1beta1 --go_opt=Mtrisa/data/generic/v1beta1/transaction.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/data/generic/v1beta1
//go:generate -command buildrpc buildproto

//go:generate buildproto --go_opt=Mtraveler/common/v1/address.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 traveler/common/v1/address.proto
//go:generate buildproto --go_opt=Mtraveler/common/v1/errors.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 traveler/common/v1/errors.proto
//go:generate buildproto --go_opt=Mtraveler/common/v1/sunrise.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 traveler/common/v1/sunrise.proto
//go:generate buildproto --go_opt=Mtraveler/common/v1/transaction.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 traveler/common/v1/transaction.proto

//go:generate buildproto --go-grpc_out=pb/ --go-grpc_opt=module=github.com/ciphertrace/apis/traveler/examples/go/pb --go_opt=Mtraveler/v1/api.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/v1 --go-grpc_opt=Mtraveler/v1/api.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/v1 --go_opt=Mtraveler/common/v1/address.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go_opt=Mtraveler/common/v1/address.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go_opt=Mtraveler/common/v1/errors.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go_opt=Mtraveler/common/v1/sunrise.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go_opt=Mtraveler/common/v1/transaction.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go-grpc_opt=Mtraveler/common/v1/address.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go-grpc_opt=Mtraveler/common/v1/address.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go-grpc_opt=Mtraveler/common/v1/errors.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go-grpc_opt=Mtraveler/common/v1/sunrise.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go-grpc_opt=Mtraveler/common/v1/transaction.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/traveler/common/v1 --go-grpc_opt=module=github.com/ciphertrace/apis/traveler/examples/go/pb --go-grpc_opt=Mivms101/enum.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/ivms101 --go-grpc_opt=Mivms101/ivms101.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/ivms101 --go-grpc_opt=Mivms101/identity.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/ivms101 --go-grpc_opt=Mtrisa/api/v1beta1/errors.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/api/v1beta1 --go-grpc_opt=Mtrisa/api/v1beta1/api.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/api/v1beta1 --go-grpc_opt=Mtrisa/data/generic/v1beta1/transaction.proto=github.com/ciphertrace/apis/traveler/examples/go/pb/trisacrypto/trisa/data/generic/v1beta1 traveler/v1/api.proto

//go:generate buildproto --proto_path=../../../trisacrypto ivms101/enum.proto ivms101/identity.proto ivms101/ivms101.proto trisa/api/v1beta1/api.proto trisa/api/v1beta1/errors.proto trisa/data/generic/v1beta1/transaction.proto
