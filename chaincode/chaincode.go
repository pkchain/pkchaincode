package main // import "gamecenter.mobi/paicode/chaincode"

import (
	"errors"
	"fmt"
	"sync"
	_ "strconv"
	_ "encoding/hex"
	
	"github.com/hyperledger/fabric/peerex"	
	"github.com/hyperledger/fabric/core/chaincode/shim"	
	
	persistpb "gamecenter.mobi/paicode/protos"
	sec "gamecenter.mobi/paicode/chaincode/security" 
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	_ "gamecenter.mobi/paicode/transactions"
)

type paiStatus struct{
	totalPai 	int64
	frozenPai 	int64
}

type PaiChaincode struct {
	globalLock 	sync.RWMutex
	cacheOK    	bool
	paistat		paiStatus
}

const (

	global_setting_entry string = tx.Global_setting_entry
	
)

var privilege_Def = map[string]string{
	tx.Admin_funcs: sec.AdminPrivilege,
	tx.Manage_funcs: sec.ManagerPrivilege,
	tx.User_funcs: sec.DelegatePrivilege}

var logger = peerex.InitLogger("chaincode")

func (s *paiStatus) init(set *persistpb.DeploySetting){
	s.totalPai = set.TotalPais
	s.frozenPai = set.UnassignedPais
}

func (s *paiStatus) set(set *persistpb.DeploySetting){
	 set.TotalPais = s.totalPai
	 set.UnassignedPais = s.frozenPai
}

func (t *PaiChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	err := t.handleInit(stub, args)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *PaiChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	if err := t.updateCache(stub); err != nil{
		return nil, err
	}
	
	rolePriv, region := sec.Helper.GetPrivilege(stub)
	
	funcGrp := function[:tx.FuncPrefix]	
	expectPriv := privilege_Def[funcGrp]
	
	//check priviledge
	if !sec.Helper.VerifyPrivilege(rolePriv, expectPriv){ 
		sec.Helper.ActiveAudit(stub, fmt.Sprintf("Call function <%s> without require priviledge: <%s vs %s>", 
				function, expectPriv, rolePriv))
		return nil, errors.New("No priviledge")
	}	
	
	var err error
	switch funcGrp{
		case tx.Admin_funcs:
		case tx.Manage_funcs:
		case tx.User_funcs:
			err = t.handleUserFuncs(stub, function, region, args)
		default:
			return nil, errors.New("Function group not exist or invokable")
	}

	return nil, err
}


func (t *PaiChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if err := t.updateCache(stub); err != nil{
		return nil, err
	}

	h, ok := tx.QueryMap[function]
	if !ok{
		return nil, errors.New(fmt.Sprint("Not a registered function:", function))
	}

	return h.Handle(stub, args)
}

func main() {
	err := shim.Start(new(PaiChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
