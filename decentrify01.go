package main

//packages used are imported, just like header files in C++
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// represents our DecentrifyChaincode object
type DecentrifyChaincode struct {
}

// represents the degree asset object
type degree struct {
	ObjectType 		string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	DegreeId        int    `json:"degreeid"`    //id of the degree
	StudentName     string `json:"studentname"` //name of the degree holder student
	InstituteName   string `json:"institutename"` //name of the degree holder institute
	Duration     	int `json:"duration"` //duration of the degree
	PassingYear     int	`json:"passingyear"` //passing year of degree
	Cgpa       		float32 `json:"cgpa"` //commulative grade point avverage of degree holder
	AllowedViews    int `json:"allowedviews"` //number of views allowed for this degree
}

// Main method
// ==============================================================================================================================
func main() {
	err := shim.Start(new(DecentrifyChaincode)) //registers the instance of DecentrifyChaincode with fabric runtime
	if err != nil {
		fmt.Printf("Error starting Decentrify Chaincode: %s", err)
	}
}

// shim.ChaincodeStubInterface - which is used to access and modify the ledger, and to make invocations between chaincodes
// Init - initializes chaincode
// ==============================================================================================================================
func (t *DecentrifyChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("init executed")
	return shim.Success(nil)
}

// Invoke - our entry point for Invocations
// ==============================================================================================================================
func (t *DecentrifyChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createDegree" { //createDegree method definition - create a new degree, store into chaincode state
		return t.createDegree(stub, args)
	} else if function == "invokeDegreeAccess" { //change number of allowed views for the degree so only particular number of times it could be viewed
		return t.invokeDegreeAccess(stub, args)
	} else if function == "viewDegree" { //read a degree from state, display  and decrement allowedViews by 1
		return t.viewDegree(stub, args)
	} else if function == "revokeAccess" { //change number of allowed views to 0 so no one could view the degree
		return t.revokeAccess(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error if unknown function
	return shim.Error("Received unknown function invocation")
}

// createDegree method definition - create a new degree, store into chaincode state
// ==============================================================================================================================
func (t *DecentrifyChaincode) createDegree(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err, err0, err1, err2, err3, err4, err5, err6 error

	// according to degree struct
	// DegreeId    StudentName  InstituteName  Duration  PassingYear  Cgpa  AllowedViews
	//     0            1             2           3          4          5        6
	//    "1",   "muhammad anas",    "BU",       "4",      "2019",    "3.42",   "0"

	//checks if number of args are less than expected args
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// input validation for empty strings
	fmt.Println("- start create degree")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}

	// according to degree struct
	// DegreeId    StudentName  InstituteName  Duration  PassingYear  Cgpa  AllowedViews
	//     0            1             2           3          4          5        6
	//    "1",   "muhammad anas",    "BU",       "4",      "2019",    "3.42",   "0"

	degreeId, err0 := strconv.Atoi(args[0]) //id of the degree
	studentName, err1 := strings.ToLower(args[1]) //name of the degree holder student
	instituteName, err2 := strings.ToLower(args[2]) //name of the degree holder institute
	duration, err3 := strconv.Atoi(args[3]) //duration of the degree
	passingYear, err4 := strconv.Atoi(args[4]) //passing year of degree
	cgpa, err5 := strconv.ParseFloat(args[5], 32) //commulative grade point average of degree holder
	allowedViews, err6 := strconv.Atoi(args[6]) //number of views allowed for this degree

	// input validation for invalid strings
	if err0 != nil {
		return shim.Error("1st argument must be a numeric string")
	}
	if err1 != nil {
		return shim.Error("2nd argument must be a text string")
	}
	if err2 != nil {
		return shim.Error("3rd argument must be a text string")
	}
	if err3 != nil {
		return shim.Error("4th argument must be a numeric string")
	}
	if err4 != nil {
		return shim.Error("5th argument must be a numeric string")
	}
	if err5 != nil {
		return shim.Error("6th argument must be a numeric string")
	}
	if err6 != nil {
		return shim.Error("7th argument must be a numeric string")
	}

	// checks if the degree already exists
	degreeAsBytes, err := stub.GetState(strconv.Itoa(degreeId))
	if err != nil {
		return shim.Error("Failed to get degree: " + err.Error())
	} else if degreeAsBytes != nil {
		fmt.Println("This degree already exists: " + degreeId)
		return shim.Error("This degree already exists: " + degreeId)
	}

	// according to degree struct
	// DegreeId    StudentName  InstituteName  Duration  PassingYear  Cgpa  AllowedViews
	//     0            1             2           3          4          5        6
	//    "1",   "muhammad anas",    "BU",       "4",      "2019",    "3.42",   "0"

	// create degree object and marshal to JSON
	// marshaling is the process of transforming the memory representation of an object to a data format suitable for storage
	objectType := "degree"
	degree := &degree{objectType, degreeId, studentName, instituteName, duration, passingYear, cgpa, allowedViews}
	degreeJSONasBytes, err := json.Marshal(degree)
	if err != nil {
		return shim.Error(err.Error())
	}

	// save degree to state
	err = stub.PutState(strconv.Itoa(degreeId), degreeJSONasBytes) //writes the state variable key of value to the state store, if the variable already exists, the value will be overwritten
	if err != nil {
		return shim.Error(err.Error())
	}

	// degree saved. Return success
	fmt.Println("- end create degree")
	return shim.Success(nil)
}

// invokeDegreeAccess method definition - change allowed views, so limited number of times degree could be viewed
// ==============================================================================================================================
func (t *DecentrifyChaincode) invokeDegreeAccess(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	// according to degree struct
	// DegreeId   AllowedViews
	//    0            1
	//   "1",         "10"

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	degreeId := strconv.Atoi(args[0])
	newAllowedViews := strconv.Atoi(args[1])
	fmt.Println("- start invokeDegreeAccess ", degreeId, newAllowedViews)

	degreeAsBytes, err := stub.GetState(strconv.Itoa(degreeId))
	if err != nil {
		return shim.Error("Failed to get degree:" + err.Error())
	} 
	else if degreeAsBytes == nil {
		return shim.Error("Degree does not exist")
	}

	updatedDegree := degree{}
	err = json.Unmarshal(degreeAsBytes, &updatedDegree) //unmarshal it also known as JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	updatedDegree.AllowedViews = (updatedDegree.AllowedViews + newAllowedViews) //change the allowed views for degree

	degreeJSONasBytes, _ := json.Marshal(updatedDegree)
	err = stub.PutState(strconv.Itoa(degreeId), degreeJSONasBytes) //rewrite the degree
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end invokeDegreeAccess (success)")
	return shim.Success(nil)
}

// viewDegree method definition - read a degree from state and decrement allowedViews by 1
// ==============================================================================================================================
func (t *DecentrifyChaincode) viewDegree(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var 
	var degreeId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID of the Degree to query")
	}

	degreeId = strconv.Atoi(args[0])
	valAsbytes, err := stub.GetState(strconv.Itoa(degreeId)) //get the degree from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + degreeId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Degree does not exist: " + degreeId + "\"}"
		return shim.Error(jsonResp)
	}

	updatedDegree := degree{}
	err = json.Unmarshal(valAsBytes, &updatedDegree) //unmarshal it also known as JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	//if number of allowed view > 0 then no one could view the degree
	else if updatedDegree.AllowedViews <= 0 {
		jsonResp = "{\"Error\":\"No rights to view degree " + degreeId + "\"}"
		return shim.Error(jsonResp)
	}
	updatedDegree.AllowedViews = (updatedDegree.AllowedViews - 1) //decrement number of allowed views for the degree on each view

	degreeJSONasBytes, _ := json.Marshal(updatedDegree)
	err = stub.PutState(strconv.Itoa(degreeId), degreeJSONasBytes) //rewrite the degree after decrementing number oof allowed views
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(valAsbytes)
}

// revokeDegreeAccess method definition - change number of allowed views to 0, so no one could  view the degree
// ==============================================================================================================================
func (t *DecentrifyChaincode) revokeDegreeAccess(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	// according to degree struct
	// DegreeId
	//    0    
	//   "1",  

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID of the Degree to query")
	}

	degreeId := strconv.Atoi(args[0])
	fmt.Println("- start invokeDegreeAccess ", degreeId)

	degreeAsBytes, err := stub.GetState(strconv.Itoa(degreeId)
	if err != nil {
		return shim.Error("Failed to get degree:" + err.Error())
	} 
	else if degreeAsBytes == nil {
		return shim.Error("Degree does not exist")
	}

	updatedDegree := degree{}
	err = json.Unmarshal(degreeAsBytes, &updatedDegree) //unmarshal it also known as JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	updatedDegree.AllowedViews = 0 //change the allowed views for degree to 0

	degreeJSONasBytes, _ := json.Marshal(updatedDegree)
	err = stub.PutState(strconv.Itoa(degreeId), degreeJSONasBytes) //rewrite the degree
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end revokeDegreeAccess (success)")
	return shim.Success(nil)
}