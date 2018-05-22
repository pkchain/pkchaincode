package transaction 

import (
	"github.com/op/go-logging"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	persistpb "gamecenter.mobi/paicode/protos"
)	

//each function name has a 4-bytes prefix
const (
	
	Global_setting_entry string = "global_setting"
	
	FuncPrefix int = 4
	
	Admin_funcs string = "ADMN"
	Manage_funcs string = "MANG"
	User_funcs string = "USER"
	Query_funcs string = "QURY"
	
)

var logger = logging.MustGetLogger("transaction")

var UserFund string = User_funcs + "_FUND"
var UserRegPublicKey string = User_funcs + "_REGPUBLICKEY"
var UserAuthChaincode string = User_funcs + "_AUTHCHAINCODE"

var QueryUser string = Query_funcs + "_USER"
var QueryNode string = Query_funcs + "_NODE"
var QueryGlobal string = Query_funcs + "_GLOBAL"
var QueryRec string = Query_funcs + "_RECORD"

type AdminOrManageTx interface{
	Handle(*persistpb.DeploySetting, shim.ChaincodeStubInterface, []string) (*persistpb.DeploySetting, map[string]*persistpb.UserData, error) 
}

var AdminMap = map[string]AdminOrManageTx{} 

type UserTx interface{
	Handle(string, *persistpb.UserData, shim.ChaincodeStubInterface, []string) (map[string]*persistpb.UserData, error) 
}

var UserTxMap = map[string]UserTx{} 

type Query interface{
	Handle(shim.ChaincodeStubInterface, []string) ([]byte, error) 
}

var QueryMap = map[string]Query{} 


