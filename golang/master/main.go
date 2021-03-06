package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type PrivateDataCC struct {
	CommonChaincode
}

const (
	name       = "PrivateDataCC"
	collection = "private1"
	counterKey = "iterator"
)

func (t *PrivateDataCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Logger.Info(" Init ")
	t.Prepare(stub)
	t.PutState(counterKey, []byte(strconv.Itoa(0)))
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *PrivateDataCC) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	//defer golang.PanicDefer(&response)
	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Invoke", fcn)
	t.Logger.Debug(fcn, "params", params)
	var responseBytes []byte
	switch fcn {
	case "putPrivate":
		var CN = t.GetThisCreator().Certificate.Subject.CommonName
		var txTime = UnixMilliSecond(t.GetTxTime()).String()
		t.PutPrivateData(collection, collection, []byte(CN+"|"+txTime))
	case "getPrivate":
		var pData = t.GetPrivateData(collection, collection)
		t.Logger.Info("pData" + string(pData))
	case "increase":

		var old = Atoi(string(t.GetState(counterKey)))
		var i = old + 1
		var iBytes = []byte(strconv.Itoa(i))
		t.PutState(counterKey, iBytes)

		responseBytes = iBytes
	case "get":
		var key = params[0]
		responseBytes = t.GetState(key)
	case "getBinding":
		//TODO what is this
		responseBytes = []byte(HexEncode(t.GetBinding()))
	case "getDecorations":
		//TODO what is this
		responseBytes = ToJson(stub.GetDecorations())
	default:

	}

	return shim.Success(responseBytes)
}

func main() {
	var cc = PrivateDataCC{}
	cc.SetLogger(name)
	ChaincodeStart(&cc)
}
