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
	"errors"
	"fmt"
	"strconv"

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

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 4{
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var party, operation, user string
	var balance int
	var err error
	party = args[0]
	operation = args[1]
	user = args[2]
	balance, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for balance")
	}
	err = stub.PutState(party, []byte(operation))
    if err != nil {
        return nil, err
    }
	err = stub.PutState(user, []byte(strconv.Itoa(balance)))
    if err != nil {
        return nil, err
    }

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init",args)
	} 
	if function == "read" {
        return t.read(stub,args)
    }
	if function == "write" {
        return t.write(stub, args)
    }
	if function == "transfer" {
        return t.transfer(stub, args)
    }
	
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("running transfer()")

    if len(args) != 5 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }
	var party, operation, user1, user2 string
	var amount, balance1, balance2 int
	var err error
    party = args[0]                           
    operation = args[1]
	user1 = args[2]
	user2 = args[3]
	amount, err = strconv.Atoi(args[4])
	if err != nil {
		return nil, errors.New("Expecting integer value for transfer")
	}
	bal1,err :=  stub.GetState(user1)
	balance1, _ = strconv.Atoi(string(bal1))
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if balance1<amount{
		return nil, errors.New("less balance")
	}
	bal2,err := stub.GetState(user2)
	balance2, _ = strconv.Atoi(string(bal2))
	balance1 = balance1-amount
	balance2 = balance2+amount
    err = stub.PutState(user1,[]byte(strconv.Itoa(balance1)))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
	err = stub.PutState(user2,[]byte(strconv.Itoa(balance2)))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
	err = stub.PutState(party,[]byte(operation))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    
    if len(args) != 4{
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var party, operation, user string
	var balance int
	var err error
	party = args[0]
	operation = args[1]
	user = args[2]
	balance, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for balance")
	}
	err = stub.PutState(party, []byte(operation))
    if err != nil {
        return nil, err
    }
	err = stub.PutState(user, []byte(strconv.Itoa(balance)))
    if err != nil {
        return nil, err
    }

	return nil, nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    
    if len(args) != 3{
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var party, operation, user string
	var err error
	party = args[0]
	operation = args[1]
	user = args[2]
	balance, err := stub.GetState(user)
	if err != nil {
		return nil, errors.New("Expecting integer value for balance")
	}
	err = stub.PutState(party, []byte(operation))
    if err != nil {
        return nil, err
    }
	err = stub.PutState(user, balance)
    if err != nil {
        return nil, err
    }

	return nil, nil
}