package main // import "gamecenter.mobi/paicode/chaincode"

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"	
	proto "github.com/golang/protobuf/proto"
	
	persistpb "gamecenter.mobi/paicode/protos"
	pb "gamecenter.mobi/paicode/protos"
	sec "gamecenter.mobi/paicode/chaincode/security" 
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	txutil "gamecenter.mobi/paicode/transactions"
)

func (t *PaiChaincode) updateCache(stub shim.ChaincodeStubInterface) error{
	t.globalLock.RLock()
	
	if !t.cacheOK{
		t.globalLock.RUnlock()
		t.globalLock.Lock()
		defer t.globalLock.Unlock()
		
		rawset, err := stub.GetState(global_setting_entry)
		if err != nil{
			return err
		}
		
		if rawset == nil{
			return errors.New("FATAL: No global setting found")
		}
		
		setting := &persistpb.DeploySetting{}
		err = proto.Unmarshal(rawset, setting)
		
		if err != nil{
			return err
		}
		
		sec.InitSecHelper(setting)
		t.paistat.init(setting)
		txutil.AddrHelper = txutil.AddressHelper(setting.NetworkCode)
		logger.Info("Update global setting:", setting)
		
		t.cacheOK = true	
	}else{
		t.globalLock.RUnlock()
	}
	
	return nil
}

func saveGlobalStatus(s *paiStatus, stub shim.ChaincodeStubInterface) error{	
	
	//a Write After Read process
	rawset, err := stub.GetState(global_setting_entry)
	if err != nil{
		return err
	}
	
	if rawset == nil{
		return errors.New("FATAL: No global setting found")
	}
	
	setting := &persistpb.DeploySetting{}
	err = proto.Unmarshal(rawset, setting)
	
	if err != nil{
		return err
	}
	
	s.set(setting)
	
	logger.Info("Save current global setting:", setting)	
	
	rawset, err = proto.Marshal(setting)
	if err != nil{
		return err
	}
	
	return stub.PutState(global_setting_entry, rawset)	
}

func (t *PaiChaincode) handleInit(stub shim.ChaincodeStubInterface, args []string) error{
	if args == nil{
		return errors.New("No init argument")	
	}
	
	settings := &pb.InitChaincode{}
	err := txutil.DecodeChaincodeTx(args[0], settings)
	if err != nil{
		return err
	}
	
	var sectmp = &sec.SecurityPolicy{}
	sectmp.Update(settings.Mainsetting)
	priv, _ := sectmp.GetPrivilege(stub)
	if !sectmp.VerifyPrivilege(priv, sec.AdminPrivilege){
		return errors.New("Not a admin for paicode")
	}
	
	var totalassigned int64 = 0
	//assign additional user info
	for _, ud := range settings.PreassignedUser{
		
		logger.Info("Preassign", ud.Pais, "pais to user", ud.Userid)
		rawset, err := proto.Marshal(&persistpb.UserData{Pais: ud.Pais})
		if err != nil{
			return err
		}
		err = stub.PutState(ud.Userid, rawset)
		if err != nil{
			return err
		}
		
		totalassigned += ud.Pais		
	}
	
	if totalassigned > settings.Mainsetting.TotalPais{
		return errors.New("Assigned too many pais than setting")
	}
	
	if settings.Mainsetting.TotalPais - totalassigned != settings.Mainsetting.UnassignedPais {
		settings.Mainsetting.UnassignedPais = settings.Mainsetting.TotalPais - totalassigned
		logger.Warning("Unmatch amount for rest pais, adjust to", settings.Mainsetting.UnassignedPais)
	}
	
	logger.Info("Init verify OK, persisting global settings")		
	rawset, err := proto.Marshal(settings.Mainsetting)
	if err != nil{
		return err
	}	
	
	err = stub.PutState(global_setting_entry, rawset)
	if err != nil{
		return err
	}
	
	logger.Info("Init chaincode for network", settings.Mainsetting.NetworkCode, 
		"finished, debugMode:", settings.Mainsetting.DebugMode)
	
	return t.updateCache(stub)
}

func (t *PaiChaincode) handleAdminManageFuncs(stub shim.ChaincodeStubInterface, function string, args []string) error{
	
	h, ok := tx.AdminMap[function]
	if !ok{
		return errors.New(fmt.Sprint("Not a registered function:", function))
	}
	
	t.globalLock.Lock()
	defer t.globalLock.Unlock()
	
	globalset := &persistpb.DeploySetting{}
	t.paistat.set(globalset)
	
	globalout, outuds, err := h.Handle(globalset, stub, args)
	
	if err != nil{
		return err
	}
	
	//update cache
	t.paistat.init(globalout)
		
	err = saveGlobalStatus(&t.paistat, stub)
	if err != nil{
		return err
	}	
	
	for id, ud := range outuds{
		raw, err := proto.Marshal(ud)
		if err != nil{
			return err
		}
		
		err = stub.PutState(id, raw)
		if err != nil{
			return err
		}
	}
	
	return nil	
}

func (t *PaiChaincode) handleUserFuncs(stub shim.ChaincodeStubInterface, function string, region string, args []string) error{
	
	h, ok := tx.UserTxMap[function]
	if !ok{
		return errors.New(fmt.Sprint("Not a registered function:", function))
	}

	cs := txutil.UserTxConsumer{}
	err := cs.ParseArgumentsFirst(args)
	if err != nil{
		return err
	}	
	
	raw, err := stub.GetState(cs.GetUserId())
	if err != nil{
		return err
	}

	userdata := &persistpb.UserData{}
	
	if raw == nil {
		return errors.New("No corresponding user")
	}
	
	err = proto.Unmarshal(raw, userdata)	
	if err != nil{
		return err
	}
	
	if function == tx.UserRegPublicKey {
		/*region is set here*/
		userdata.ManagedRegion = region		
	}else{
		if userdata.Pk == nil{
			return errors.New("User has not register his public key yet")
		}
		
		/*check wether region is mathced*/
		if !sec.Helper.VerifyRegion(userdata.ManagedRegion, region){
			return errors.New(fmt.Sprint("User tx is invoked in different region:", region))
		}
	}
	
	outuds, err := h.Handle(cs.GetUserId(), userdata, stub, args)
	if err != nil{
		return err
	}
	
	for id, ud := range outuds{
		raw, err := proto.Marshal(ud)
		if err != nil{
			return err
		}
		
		err = stub.PutState(id, raw)
		if err != nil{
			return err
		}
	}
	
	return nil
}
