package client

import (
	
	"fmt"
	"strings"
	
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	txutil "gamecenter.mobi/paicode/transactions"
	pb "gamecenter.mobi/paicode/protos"
	
	"github.com/hyperledger/fabric/peerex"	
)

type blockExplorer struct{
	
}

type fundTxSimple struct{
	From   string `json:"from"`
	To     string `json:"to"`
	Message []byte `json:"comment"`
	Amount  int   `json:"amount"`
}

//decode payload: <payload>
func (_* blockExplorer) DecodePayload(args ...string) (*fundTxSimple, error){
	if len(args) != 1{
		return nil, fmt.Errorf("Require payload string")
	}
	
	invoke, err := peerex.DecodeTransactionToInvoke(args[0])
	
	if err != nil{
		return nil, fmt.Errorf("Decode payload fail", err.Error())
	}
	
	fund := string(invoke.ChaincodeSpec.CtorMsg.Args[0])
	
	if strings.Compare(fund, tx.UserFund) != 0{
		return nil, fmt.Errorf("Wrong function:", fund)
	}
	
	cm := txutil.UserTxConsumer{}
	v1 := &pb.Fund{}
	v2 := &pb.Funddata{}
	iargs := make([]string, len(invoke.ChaincodeSpec.CtorMsg.Args[1:]))
	
	for i, a := range invoke.ChaincodeSpec.CtorMsg.Args[1:]{
		iargs[i] = string(a)
	}
	
	if len(iargs) < 4{
		err = cm.ParseArguments(iargs, v1)	
	}else{
		err = cm.ParseArguments(iargs, v1, v2)	
	}
	
	if err != nil{
		return nil, fmt.Errorf("Decode chaincode arguments fail", err.Error())
	}	
	
	var vdata *pb.Funddata
	if v1.InvokeChaincode != 0{
		vdata = v2
	}else{
		vdata = v1.GetUserfund()
	}
	
	if vdata == nil{
		return nil, fmt.Errorf("No valid fundtx data")
	}
	
	return &fundTxSimple{cm.GetUserId(), vdata.ToUserId, cm.GetTxNounce(), int(vdata.Pai)}, nil
		
}