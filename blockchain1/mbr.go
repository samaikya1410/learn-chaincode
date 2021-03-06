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

type Entity struct{
	Entity_Name string `json:"entity_name"`
	Entity_Role string `json:"entity_role"`
}

type Transaction struct{
	Entity_Involved Entity `json:"entity_involved"`
	Operation string `json:"operation"`
	Claim_Id string `json:"claim_id"`
	Bill_Id string `json:"bill_id"`
	Bill_Details string `json:"bill_details"`
	Bill_Status string `json:"bill_status"`
	Date string `json:"date"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0{
		return nil, errors.New("Incorrect number of arguments while deploying.")
	}
	var e Entity
	var c int
	var err error
	c = 0
	err = stub.PutState("count",[]byte(strconv.Itoa(c)))
	if err != nil {
				return nil, err
	}
	e.Entity_Name = "admin"
	e.Entity_Role = "admin"
	ebytes, err := json.Marshal(&e)
	if err != nil {
		return nil, errors.New("Error deploying chaincode")
	}
	err = stub.PutState("admin",ebytes)
	if err != nil {
        return nil, err
  }
	return nil, nil
}


func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3{
		return nil, errors.New("Incorrect number of arguments while registering.")
	}
	var e Entity
	e.Entity_Name = args[1]
	e.Entity_Role = args[2]
	ebytes, err := json.Marshal(&e)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(args[0],ebytes)
	if err != nil {
        return nil, err
  }
	return nil, nil
}


func (t *SimpleChaincode) transact(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 6{
		return nil, errors.New("Incorrect number of arguments while transacting.")
	}
	var tr Transaction
	var user Entity
	var err error
	tr.Entity_Involved = user
	tr.Operation = args[1]
	tr.Claim_Id = args[2]
	tr.Bill_Id = args[3]
	tr.Bill_Details = args[4]
	tr.Bill_Status = args[5]
	tr.Date = time.Now().String()

	var c int
	var cstring string
	//get current transaction count, increment it
	cbytes, err := stub.GetState("count")
	if err != nil {
			 return nil, err
	}
	c, _ = strconv.Atoi(string(cbytes))
	c = c+1
	cstring = strconv.Itoa(c)
	err = stub.PutState("count",[]byte(cstring))
	if err != nil {
			 return nil, err
	}
	//enter the Transaction
	userbytes, err := stub.GetState(args[0])
	if err != nil {
			 return nil, errors.New("username wrong.")
	}
	err = json.Unmarshal(userbytes, &user)
	if err != nil {
			 return nil, err
	}
	tbytes, err := json.Marshal(&tr)
	if err != nil {
			 return nil, err
	}
	err = stub.PutState(cstring,tbytes)
	if err != nil {
			 return nil, err
	}
	return nil,nil
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "register" {													//initialize the chaincode state, used as reset
		return t.register(stub,args)
	}
	if function == "transact" {
        return t.transact(stub, args)
  }

	fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "get_count"{
		c, err := stub.GetState("count")
		if err != nil {
	 		return nil, errors.New("cannot know the count")
	 	}
		return c, nil
	}
	if function == "get_transaction"{
		transaction, err := stub.GetState(args[0])
		if err != nil {
			return nil, fmt.Errorf("Failed getting transaction, [%v]", err)
		}
		return transaction, nil
	}
	if function == "get_user"{
		c, err := stub.GetState(args[0])
		if err != nil {
	 		return nil, errors.New("cannot get user")
	 	}
		return c, nil
	}

	return nil, errors.New("Received unknown function" )
}
