package transaction

import (
	"errors"
	"strings"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
	txutil "gamecenter.mobi/paicode/transactions"
	pb 	   "gamecenter.mobi/paicode/protos"
	persistpb "gamecenter.mobi/paicode/protos"
)

type regPublicKeyHandler struct{
	
}

type authChaicodeHandler struct{
	
}

func init(){
	UserTxMap[UserRegPublicKey] = &regPublicKeyHandler{}
}

func (f *regPublicKeyHandler) Handle(uid string, ud *persistpb.UserData, stub shim.ChaincodeStubInterface, 
	args []string) (outud map[string]*persistpb.UserData, err error) {
	
	cs := txutil.UserTxConsumer{PublicKey: nil}
	v := &pb.RegPublicKey{}
	
	//parse with no public key (not verify signature)
	err = cs.ParseArguments(args, v)
	if err != nil{
		return
	}
	
	if len(cs.GetTxNounce()) < txutil.SecurityNounceLen{
		err = errors.New("Register publickey tx require a longer nounce")
		return 
	}
	
	pk := (*txutil.PublicKey) (v.Pk).ECDSAPublicKey()
	
	//verify public key
	if strings.Compare(txutil.AddrHelper.GetUserId(pk), uid) != 0{
		logger.Warning("Userid", uid, "is invoked by an unknown publickey", ud.Pk)
		err = errors.New("Public key is not match with user id")
		return
	}
	
	//need to verify signature again ...
	cs.Reset()
	cs.PublicKey = pk
	err = cs.ParseArguments(args, v)	
	if err != nil{
		logger.Warning("Userid", uid, "get regpublic tx invoked by wrong signature")
		return
	}
	
	//do we register before? make a log
	if ud.Pk == nil{
		logger.Info(uid, "set public key at region", ud.ManagedRegion)
	}else{
		logger.Info(uid, "reset user's region to", ud.ManagedRegion)
	}
	
	ud.Pk = v.Pk
	ud.LastActive = acquireTsNow(stub)	
	outud = map[string]*persistpb.UserData{uid: ud}
	
	return
}
