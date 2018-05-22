package transaction

import (
	"errors"
	"encoding/json"
	"strings"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	proto "github.com/golang/protobuf/proto"
	sec "gamecenter.mobi/paicode/chaincode/security" 
	_ "gamecenter.mobi/paicode/transactions"
	persistpb "gamecenter.mobi/paicode/protos"
)

type queryUserHandler struct{
	
}

type queryGlobalHandler struct{
	
}

type queryNodeHandler struct{
	
}

type queryRecordHandler struct{
	
}

func init(){
	QueryMap[QueryUser] = &queryUserHandler{}
	QueryMap[QueryGlobal] = &queryGlobalHandler{}
	QueryMap[QueryNode] = &queryNodeHandler{}
	QueryMap[QueryRec] = &queryRecordHandler{}
}

func (_ *queryUserHandler) Handle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
	if len(args) != 1{
		return nil, errors.New("Invalid query argument")
	}

	//TODO: decode some additional options ...
	//if len(args) > 1{
	//	err := txutil.DecodeChaincodeTx(args[1], pb.QueryUser)
	//	if err != nil{
	//		return nil, err
	//	}		
	//}
	
	raw, err := stub.GetState(args[0])
	if err != nil{
		return nil, err
	}

	userdata := &persistpb.UserData{}
	
	if raw == nil {
		return nil, errors.New("No corresponding user")
	}
	
	err = proto.Unmarshal(raw, userdata)	
	if err != nil{
		return nil, err
	}
	
	return json.Marshal(userdata)
} 

func (_ *queryGlobalHandler) Handle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
	if len(args) != 0{
		return nil, errors.New("Invalid query argument")
	}	
	
	rawset, err := stub.GetState(Global_setting_entry)
	if err != nil{
		return nil, err
	}
	
	if rawset == nil{
		return nil, errors.New("No global setting found")
	}
	
	setting := &persistpb.DeploySetting{}
	err = proto.Unmarshal(rawset, setting)
	
	if err != nil{
		return nil, err
	}
	
	return json.Marshal(setting)
} 

func (_ *queryRecordHandler) Handle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
	if len(args) != 1{
		return nil, errors.New("Invalid query argument")
	}
	
	if strings.Compare(FundNouncePrefix, args[0][:3]) != 0{
		return nil, errors.New("Not a fundtx record key")
	}
	
	rawset, err := stub.GetState(args[0])
	if err != nil{
		return nil, err
	}
	
	if rawset == nil{
		return nil, errors.New("fundtx record not found")
	}
	
	record := &persistpb.NounceData{}
	err = proto.Unmarshal(rawset, record)
	
	if err != nil{
		return nil, err
	}
	
	return json.Marshal(record)	
}

type nodeInfo struct{
	Region string `json:"region"`
	Priv string `json:"privilege"`
}

func (_ *queryNodeHandler) Handle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	rolePriv, region := sec.Helper.GetPrivilege(stub)
	
	return json.Marshal(nodeInfo{Priv: rolePriv, Region: region})
}

