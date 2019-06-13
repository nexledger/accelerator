/*
 *    Copyright 2019 Samsung SDS
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type PingPongChaincode struct {
}

func (t *PingPongChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *PingPongChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fnc := string(stub.GetArgs()[0])
	switch fnc {
	case "ping":
		return Invoke(stub, t.ping)
	case "pong":
		return Invoke(stub, t.pong)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'insert', 'query'")
}

func (t *PingPongChaincode) ping(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := stub.PutState(args[0], []byte(args[1])); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(nil)
	}
}

func (t *PingPongChaincode) pong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if value, err := stub.GetState(args[0]); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(value)
	}
}

func main() {
	err := shim.Start(new(PingPongChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %v \n", err)
	}
}
