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
	Claim_Id string `json:"claim_id"`
	Bill_Id string `json:"bill_id"`
	Operation string `json:"operation"`
	Bill_Details string `json:"bill_details"`
	Bill_Status string `json:"bill_status"`
	Date string `json:"date"`
}
type Approve_bill struct{
	Provider_Name string `json:"provider_name"`
	Provider_Role string `json:"provider_role"`
	Claim_Id string `json:"claim_id"`
	Bill_Id string `json:"bill_id"`
	Operation string `json:"operation"`
	Bill_Details string `json:"bill_details"`
	Bill_Status string `json:"bill_status"`
	Date string `json:"date"`
}
type Pay_bill struct{
	Vendor_Name string `json:"vendor_name"`
	Provider_Name string `json:"provider_name"`
	Provider_Role string `json:"provider_role"`
	Claim_Id string `json:"claim_id"`
	Bill_Id string `json:"bill_id"`
	Approver_Name string `json:"approver_name"`
	Approver_Role string `json:"approver_role"`
	Operation string `json:"operation"`
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
	err = stub.CreateTable("to_be_validated_bills", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Provider_Name", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Provider_Role", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Claim_Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Bill_Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Operation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Bill_Details", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Bill_Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Date", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating approved_bills table")
	}
	err = stub.CreateTable("to_be_paid_bills", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Vendor", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Provider_Name", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Provider_Role", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Claim_Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Bill_Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Approver_Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Approver_Role", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Operation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Bill_Details", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Bill_Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Date", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating paid_bills table")
	}
	return nil, nil
}


func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	if e.Entity_Role == "Provider"{
		err = stub.PutState(e.Entity_Name, []byte(args[3]))
		if err != nil {
	        return nil, err
	  }
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
	userbytes, err := stub.GetState(args[0])
	if err != nil {
			 return nil, errors.New("username wrong.")
	}
	err = json.Unmarshal(userbytes, &user)
	if err != nil {
			 return nil, err
	}
	tr.Entity_Involved = user
	tr.Claim_Id = args[1]
	tr.Bill_Id = args[2]
	tr.Operation = args[3]
	tr.Bill_Details = args[4]
	tr.Bill_Status = args[5]
	tr.Date = time.Now().String()

	if(tr.Operation == "Submit"){
		bool, err := stub.InsertRow("to_be_validated_bills", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: user.Entity_Name}},
				&shim.Column{Value: &shim.Column_String_{String_: user.Entity_Role}},
				&shim.Column{Value: &shim.Column_String_{String_: tr.Claim_Id}},
				&shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Id}},
				&shim.Column{Value: &shim.Column_String_{String_: tr.Operation}},
				&shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Details}},
				&shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Status}},
				&shim.Column{Value: &shim.Column_String_{String_: tr.Date}},
			},
		})
		if (!bool && err == nil){
			return nil, errors.New("already submited")
		}
		if (!bool && err != nil){
			return  nil, errors.New("could not insert row in to_be_validated_bills")
		}
	}

	if(tr.Operation == "Validate"){
		var provider_name, provider_role string
		provider_name = args[6]
		provider_role = args[7]
		if(tr.Bill_Status == "Approved"){
			vendorbytes, err := stub.GetState(args[6])
			if err != nil {
					 return nil, errors.New("cannot get vendor")
			}
			var vendor string
			vendor = string(vendorbytes)
			bool, err := stub.InsertRow("to_be_paid_bills", shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: vendor}},
					&shim.Column{Value: &shim.Column_String_{String_: provider_name}},
					&shim.Column{Value: &shim.Column_String_{String_: provider_role}},
					&shim.Column{Value: &shim.Column_String_{String_: tr.Claim_Id}},
					&shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Id}},
					&shim.Column{Value: &shim.Column_String_{String_: user.Entity_Name}},
					&shim.Column{Value: &shim.Column_String_{String_: user.Entity_Role}},
					&shim.Column{Value: &shim.Column_String_{String_: tr.Operation}},
					&shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Details}},
					&shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Status}},
					&shim.Column{Value: &shim.Column_String_{String_: tr.Date}},
				},
			})
			if (!bool && err == nil){
				return  nil, errors.New("already approved")
			}
			if (!bool && err != nil){
				return  nil, err
			}
		}
		err = stub.DeleteRow("to_be_validated_bills", []shim.Column{
			shim.Column{Value: &shim.Column_String_{String_: provider_name}},
			shim.Column{Value: &shim.Column_String_{String_: provider_role}},
			shim.Column{Value: &shim.Column_String_{String_: tr.Claim_Id}},
			shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Id}},
		},
		)
		if err != nil {
		return nil, err
		}
	}

	if(tr.Operation == "Pay"){
		var provider_name, provider_role string
		provider_name = args[6]
		provider_role = args[7]
		err = stub.DeleteRow("to_be_paid_bills", []shim.Column{
			shim.Column{Value: &shim.Column_String_{String_: user.Entity_Name}},
			shim.Column{Value: &shim.Column_String_{String_: provider_name}},
			shim.Column{Value: &shim.Column_String_{String_: provider_role}},
			shim.Column{Value: &shim.Column_String_{String_: tr.Claim_Id}},
			shim.Column{Value: &shim.Column_String_{String_: tr.Bill_Id}},
		},
		)
		if err != nil {
		return nil, err
		}
	}

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

	if function == "get_to_be_validated_bills"{
		var rowChannel  <-chan shim.Row
		rowChannel, err := stub.GetRows("to_be_validated_bills", []shim.Column{})
		if err != nil {
	 		return nil, errors.New("cannot get rows")
	 	}
		//var rowstrings []string
		var list_of_approve_bills []Approve_bill
		var approve_bill Approve_bill
		for {
			select {
			case row, ok := <-rowChannel:
				if !ok {
					rowChannel = nil
				} else {
					//rowstrings = append(rowstrings, row.String())
					approve_bill.Provider_Name = row.Columns[0].GetString_()
					approve_bill.Provider_Role = row.Columns[1].GetString_()
					approve_bill.Claim_Id = row.Columns[2].GetString_()
					approve_bill.Bill_Id = row.Columns[3].GetString_()
					approve_bill.Operation = row.Columns[4].GetString_()
					approve_bill.Bill_Details = row.Columns[5].GetString_()
					approve_bill.Bill_Status = row.Columns[6].GetString_()
					approve_bill.Date = row.Columns[6].GetString_()
					list_of_approve_bills = append(list_of_approve_bills, approve_bill)
				}
			}
			if rowChannel == nil{
				break
			}
		}
		approve_bill_bytes, err := json.Marshal(list_of_approve_bills)
		if err != nil {
			return nil, fmt.Errorf("rows operation failed. Error marshaling JSON: %s", err)
		}
		return approve_bill_bytes, nil
	}

	if function == "get_to_be_paid_bills"{
		var vendor = args[0]
		var rowChannel  <-chan shim.Row
		rowChannel, err := stub.GetRows("to_be_paid_bills", []shim.Column{
			shim.Column{Value: &shim.Column_String_{String_: vendor}},
		})
		if err != nil {
	 		return nil, errors.New("cannot get rows")
	 	}
		//var rowstrings []string
		var list_of_pay_bills []Pay_bill
		var pay_bill Pay_bill
		for {
			select {
			case row, ok := <-rowChannel:
				if !ok {
					rowChannel = nil
				} else {
					//rowstrings = append(rowstrings, row.String())
					pay_bill.Vendor_Name = row.Columns[0].GetString_()
					pay_bill.Provider_Name = row.Columns[1].GetString_()
					pay_bill.Provider_Role = row.Columns[2].GetString_()
					pay_bill.Claim_Id = row.Columns[3].GetString_()
					pay_bill.Bill_Id = row.Columns[4].GetString_()
					pay_bill.Approver_Name = row.Columns[5].GetString_()
					pay_bill.Approver_Role = row.Columns[6].GetString_()
					pay_bill.Operation = row.Columns[7].GetString_()
					pay_bill.Bill_Details = row.Columns[8].GetString_()
					pay_bill.Bill_Status = row.Columns[9].GetString_()
					pay_bill.Date = row.Columns[10].GetString_()

					list_of_pay_bills = append(list_of_pay_bills, pay_bill)
				}
			}
			if rowChannel == nil{
				break
			}
		}
		pay_bill_bytes, err := json.Marshal(list_of_pay_bills)
		if err != nil {
			return nil, fmt.Errorf("rows operation failed. Error marshaling JSON: %s", err)
		}
		return pay_bill_bytes, nil
	}
	return nil, errors.New("Received unknown function" )
}
