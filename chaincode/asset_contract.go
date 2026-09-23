/*
 * IT Asset Management Chaincode
 *
 * Hyperledger Fabric smart contract for managing IT asset
 * lifecycle with immutable audit trail.
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AssetChaincode struct{}

type Asset struct {
	AssetID      string    `json:"assetId"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	SerialNumber string    `json:"serialNumber"`
	Owner        string    `json:"owner"`
	Department   string    `json:"department"`
	Location     string    `json:"location"`
	Status       string    `json:"status"`
	PurchaseDate string    `json:"purchaseDate"`
	Value        float64   `json:"value"`
	LastModified string    `json:"lastModified"`
	History      []Event   `json:"history"`
}

type Event struct {
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Timestamp string `json:"timestamp"`
	Details   string `json:"details"`
}

func (t *AssetChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *AssetChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "registerAsset":
		return t.registerAsset(stub, args)
	case "transferAsset":
		return t.transferAsset(stub, args)
	case "updateStatus":
		return t.updateStatus(stub, args)
	case "queryAsset":
		return t.queryAsset(stub, args)
	case "queryByOwner":
		return t.queryByOwner(stub, args)
	case "getHistory":
		return t.getHistory(stub, args)
	case "decommission":
		return t.decommission(stub, args)
	default:
		return shim.Error("Invalid function: " + function)
	}
}

func (t *AssetChaincode) registerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Expected JSON asset payload")
	}

	var asset Asset
	err := json.Unmarshal([]byte(args[0]), &asset)
	if err != nil {
		return shim.Error("Invalid asset JSON: " + err.Error())
	}

	existing, _ := stub.GetState(asset.AssetID)
	if existing != nil {
		return shim.Error("Asset already exists: " + asset.AssetID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	asset.Status = "active"
	asset.LastModified = now
	asset.History = []Event{{
		Action:    "registered",
		Actor:     asset.Owner,
		Timestamp: now,
		Details:   "Asset registered in system",
	}}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return shim.Error("JSON marshal error: " + err.Error())
	}

	err = stub.PutState(asset.AssetID, assetJSON)
	if err != nil {
		return shim.Error("Failed to register asset: " + err.Error())
	}

	return shim.Success(assetJSON)
}

func (t *AssetChaincode) transferAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Expected: assetId, newOwner, newDepartment")
	}

	assetID := args[0]
	newOwner := args[1]
	newDepartment := args[2]

	assetJSON, err := stub.GetState(assetID)
	if err != nil || assetJSON == nil {
		return shim.Error("Asset not found: " + assetID)
	}

	var asset Asset
	json.Unmarshal(assetJSON, &asset)

	now := time.Now().UTC().Format(time.RFC3339)
	event := Event{
		Action:    "transferred",
		Actor:     newOwner,
		Timestamp: now,
		Details:   fmt.Sprintf("Transferred from %s (%s) to %s (%s)", asset.Owner, asset.Department, newOwner, newDepartment),
	}

	asset.Owner = newOwner
	asset.Department = newDepartment
	asset.LastModified = now
	asset.History = append(asset.History, event)

	updatedJSON, _ := json.Marshal(asset)
	stub.PutState(assetID, updatedJSON)

	return shim.Success(updatedJSON)
}

func (t *AssetChaincode) updateStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Expected: assetId, newStatus, actor")
	}

	assetJSON, err := stub.GetState(args[0])
	if err != nil || assetJSON == nil {
		return shim.Error("Asset not found")
	}

	var asset Asset
	json.Unmarshal(assetJSON, &asset)

	now := time.Now().UTC().Format(time.RFC3339)
	asset.Status = args[1]
	asset.LastModified = now
	asset.History = append(asset.History, Event{
		Action:    "status_change",
		Actor:     args[2],
		Timestamp: now,
		Details:   "Status changed to " + args[1],
	})

	updatedJSON, _ := json.Marshal(asset)
	stub.PutState(args[0], updatedJSON)
	return shim.Success(updatedJSON)
}

func (t *AssetChaincode) queryAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Expected: assetId")
	}
	assetJSON, err := stub.GetState(args[0])
	if err != nil || assetJSON == nil {
		return shim.Error("Asset not found")
	}
	return shim.Success(assetJSON)
}

func (t *AssetChaincode) queryByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Expected: owner")
	}
	query := fmt.Sprintf(`{"selector":{"owner":"%s"}}`, args[0])
	iter, err := stub.GetQueryResult(query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var results []Asset
	for iter.HasNext() {
		item, _ := iter.Next()
		var asset Asset
		json.Unmarshal(item.Value, &asset)
		results = append(results, asset)
	}

	resultsJSON, _ := json.Marshal(results)
	return shim.Success(resultsJSON)
}

func (t *AssetChaincode) getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Expected: assetId")
	}
	iter, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()

	var history []map[string]interface{}
	for iter.HasNext() {
		item, _ := iter.Next()
		entry := map[string]interface{}{
			"txId":      item.TxId,
			"timestamp": item.Timestamp,
			"isDelete":  item.IsDelete,
		}
		if !item.IsDelete {
			var asset Asset
			json.Unmarshal(item.Value, &asset)
			entry["value"] = asset
		}
		history = append(history, entry)
	}

	historyJSON, _ := json.Marshal(history)
	return shim.Success(historyJSON)
}

func (t *AssetChaincode) decommission(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Expected: assetId, actor")
	}
	return t.updateStatus(stub, []string{args[0], "decommissioned", args[1]})
}

func main() {
	err := shim.Start(new(AssetChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err)
	}
}
