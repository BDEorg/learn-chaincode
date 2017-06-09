package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("basicgatewaylogger")

type BasicChaincode struct {
}

// ============================================================================================================================

//custom data models
type PersonalInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	DOB       string `json:"DOB"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
}

type FinancialInfo struct {
	MonthlySalary      int `json:"monthlySalary"`
	MonthlyRent        int `json:"monthlyRent"`
	OtherExpenditure   int `json:"otherExpenditure"`
	MonthlyLoanPayment int `json:"monthlyLoanPayment"`
}

type LoanApplication struct {
	ID                     string        `json:"id"`
	PropertyId             string        `json:"propertyId"`
	LandId                 string        `json:"landId"`
	PermitId               string        `json:"permitId"`
	BuyerId                string        `json:"buyerId"`
	AppraisalApplicationId string        `json:"appraiserApplicationId"`
	SalesContractId        string        `json:"salesContractId"`
	PersonalInfo           PersonalInfo  `json:"personalInfo"`
	FinancialInfo          FinancialInfo `json:"financialInfo"`
	Status                 string        `json:"status"`
	RequestedAmount        int           `json:"requestedAmount"`
	FairMarketValue        int           `json:"fairMarketValue"`
	ApprovedAmount         int           `json:"approvedAmount"`
	ReviewerId             string        `json:"reviewerId"`
	LastModifiedDate       string        `json:"lastModifiedDate"`
}

type customEvent struct {
	Type       string `json:"type"`
	Decription string `json:"description"`
}

//=======================
func GetLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering GetLoanApplication")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing loan application ID")
	}

	var loanApplicationId = args[0]
	bytes, err := stub.GetState(loanApplicationId)
	if err != nil {
		logger.Error("Could not fetch loan application with id "+loanApplicationId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreateLoanApplication")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

	var loanApplicationId = args[0]
	var loanApplicationInput = args[1]

	err := stub.PutState(loanApplicationId, []byte(loanApplicationInput))
	if err != nil {
		logger.Error("Could not save loan application to ledger", err)
		return nil, err
	}

	var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved loan application")
	return nil, nil

}

/**
Updates the status of the loan application
**/
func UpdateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering UpdateLoanApplication")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application update")
	}

	var loanApplicationId = args[0]
	var status = args[1]

	laBytes, err := stub.GetState(loanApplicationId)
	if err != nil {
		logger.Error("Could not fetch loan application from ledger", err)
		return nil, err
	}
	var loanApplication LoanApplication
	err = json.Unmarshal(laBytes, &loanApplication)
	loanApplication.Status = status

	laBytes, err = json.Marshal(&loanApplication)
	if err != nil {
		logger.Error("Could not marshal loan application post update", err)
		return nil, err
	}

	err = stub.PutState(loanApplicationId, laBytes)
	if err != nil {
		logger.Error("Could not save loan application post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'loanApplicationUpdate', description:" + loanApplicationId + "' Successfully updated status'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated loan application")
	return nil, nil

}


func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
	logger.Debug("Entering GetCertAttribute")
	attr, err := stub.ReadCertAttribute(attributeName)
	if err != nil {
		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
	}
	attrString := string(attr)
	return attrString, nil
}

//Write (function write)
func (t *BasicChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var name, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
    }

    name = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(name, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}

   

//Read (function read)
func (t *BasicChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    }

    name = args[0]
    valAsbytes, err := stub.GetState(name)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}


// main
func main() {

	// enable debug log mode 
	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)
	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	// basic startup 
	err := shim.Start(new(BasicChaincode))
	if err != nil {
		logger.Error("Could not start BasicChaincode")
	} else {
		logger.Info("BasicChaincode successfully started")
	}

}

// Init resets all the things
func (t *BasicChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    //Coding here like : 
	//err := stub.PutState("hello_world", []byte(args[0]))
    //if err != nil {
    //    return nil, err
    //}
	return nil, nil
}

// Query is our entry point for queries
func (t *BasicChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    //Coding here like : 
	 if function == "GetLoanApplication" {
	 	return GetLoanApplication(stub, args)			  // call UDF GetLoanApplication	
	 }
	//
    //if function == "read" {                            // read a variable
    //    return t.read(stub, args)                      // call UDF read
    //}	 

    return nil,  nil 
	//errors.New("Received unknown function query")
}

// Invoke is our entry point to invoke a chaincode function
func (t *BasicChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    //Coding here like : 	
	 if function == "CreateLoanApplication" {					// param compare
	 	username, _ := GetCertAttribute(stub, "username")		// call UDF GetCertAttribute
	 	role, _ := GetCertAttribute(stub, "role") 				// call UDF GetCertAttribute
	 	if role == "Bank_Home_Loan_Admin" {						// var compare
	 		return CreateLoanApplication(stub, args)			// call UDF CreateLoanApplication
	 	} else {
	 		return nil, errors.New(username + " with role " + role + " does not have access to create a loan application")
	 	}
	//
    //if function == "init" {
    //    return t.Init(stub, "init", args)
    //} else if function == "write" {
    //    return t.write(stub, args)
    //}
	

	}
	return nil, nil
}




