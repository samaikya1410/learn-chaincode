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
	Transfer_amount int `json:"transfer_amount"`
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
	var u User
	var err error
	var count, balance int
	count = 0
	balance, _= strconv.Atoi(args[3])
	u = User{Username:args[2], Balance:balance}
	reg = Register{Party:args[0], Operation:args[1], Account:u}
	regbytes, err := json.Marshal(&reg)
	if err != nil {
		fmt.Println("error creating account" + reg.Account.Username)
		return nil, errors.New("Error creating account " + reg.Account.Username)
	}
	err = stub.PutState(reg.Account.Username,regbytes)
	if err != nil {
        return nil, err
    }
	err = stub.PutState(reg.Account.Username + strconv.Itoa(count),regbytes)
	if err != nil {
        return nil, err
    }
	count = count+1
	err = stub.PutState(reg.Account.Username+"trans_count",[]byte(strconv.Itoa(count)))
	if err != nil {
        return nil, err
    }
	return nil, nil
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "register" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init",args)
	} 
	if function == "read" {
        return t.read(stub,args)
    }
	if function == "transfer" {
        return t.transfer(stub, args)
    }
	
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	var err error
	var u string
	var current int
	if function != "query" {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	u = args[0]
	u_count, err := stub.GetState(u+"trans_count")
	if err != nil {
        return nil, err
    }
	current, _ = strconv.Atoi(string(u_count))
	current = current-1
	trans, err := stub.GetState(u+strconv.Itoa(current)) 
	if err != nil {
        return trans, nil
    }

	// Handle different functions
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	

	return nil, nil
}
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("running transfer()")

    if len(args) != 5 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }
	var trans Transfer
	var user1, user2 Register
	var err error
	var count1, count2 int
    trans.Party = args[0]                           
    trans.Operation = args[1]
	trans.Sender.Username = args[2]
	trans.Reciever.Username = args[3]
	trans.Transfer_amount, err = strconv.Atoi(args[4])
	if err != nil {
		return nil, errors.New("Expecting integer value for transfer")
	}
	user1_bytes,err :=  stub.GetState(trans.Sender.Username)
	err = json.Unmarshal(user1_bytes, &user1)
	if err != nil {
		fmt.Println("Error unmarshalling user1")
		return nil, errors.New("Error unmarshalling user1")
	}
	if user1.Account.Balance < trans.Transfer_amount{
		return nil, errors.New("less balance")
	}
	user2_bytes,err := stub.GetState(trans.Reciever.Username)
	err = json.Unmarshal(user2_bytes, &user2)
	if err != nil {
		fmt.Println("Error unmarshalling user2")
		return nil, errors.New("Error unmarshalling user2")
	}
	user1.Account.Balance = user1.Account.Balance-trans.Transfer_amount
	user2.Account.Balance = user2.Account.Balance+trans.Transfer_amount
	trans.Sender.Balance = user1.Account.Balance
	trans.Reciever.Balance = user2.Account.Balance
	reg1bytes, err := json.Marshal(&user1)
    if err != nil {
        return nil, err
    }
	err = stub.PutState(user1.Account.Username,reg1bytes)  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
	reg2bytes, err := json.Marshal(&user2)
    if err != nil {
        return nil, err
    }
	err = stub.PutState(user2.Account.Username,reg2bytes)  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
	trans_bytes, err := json.Marshal(&trans)
    if err != nil {
        return nil, err
    }
	c1,err := stub.GetState(user1.Account.Username+"trans_count")
	err = stub.PutState(user1.Account.Username+string(c1),trans_bytes)  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
	c2,err := stub.GetState(user2.Account.Username+"trans_count")
	err = stub.PutState(user2.Account.Username+string(c2),trans_bytes)  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
	count1, _ = strconv.Atoi(string(c1))
	count1 = count1+1
	err = stub.PutState(user1.Account.Username+"trans_count",[]byte(strconv.Itoa(count1)))
	if err != nil {
        return nil, err
    }
	count2, _ = strconv.Atoi(string(c2))
	count2 = count2+1
	err = stub.PutState(user2.Account.Username+"trans_count",[]byte(strconv.Itoa(count2)))
	if err != nil {
        return nil, err
    }
    return nil, nil
}