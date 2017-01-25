/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

type User struct{
	Username string `json:"username"`
	Balance int `json:"balance"`
}
type Register struct{
	Party string `json:"party"`
	Operation string `json:"operation"` 
	Account User `json:"account"`
}
type Transfer struct{
	Party string `json:"party"`
	Operation string `json:"operation"`
	Sender User `json:"sender"`
	Reciever User `json:"reciever"`
}

type Read struct{
	Party string `json:"party"`
	Operation string `json:"operation"` 
	Account User `json:"account"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 4{
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var reg Register
	var err error
	var count int
	count = 0
	reg = Reciever{Party:args[0], Operation:args[1], Account{Username:args[2], Balance:strconv.Atoi(args[3])}}
	regbytes, err := json.Marshal(&reg)
	if err != nil {
		fmt.Println("error creating account" + reg.Account.User)
		return nil, errors.New("Error creating account " + reg.Account.User)
	}
	err = stub.PutState(reg.Account.Username+strconv.Itoa(count),regbytes)
	if err != nil {
        return nil, err
    }
	count=1
	err = stub.PutState("trans_count",count)
	if err != nil {
        return nil, err
    }
	return nil, nil
}
