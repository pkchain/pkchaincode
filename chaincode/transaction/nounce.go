package transaction

import (
	"errors"
	"encoding/base64"
	"crypto/sha256"
	
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	txutil "gamecenter.mobi/paicode/transactions"
	pb 	   "gamecenter.mobi/paicode/protos"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
 
const (
	FundNouncePrefix = "FNC"
	
	nounce_reuse_interval_sec int64 = 2592000 //30 days
)


type NounceManager struct{
	Tsnow *timestamp.Timestamp
	nouncekey	[]byte
	stub 		shim.ChaincodeStubInterface
}

func GenFuncNounceKeyStr(nouncekey []byte) string{
	return FundNouncePrefix + base64.StdEncoding.EncodeToString(nouncekey)
}

func GenfundNounce(from string, to string, nounce []byte) []byte {
	idfrombyte := txutil.AddrHelper.DecodeUserid(from)
	if idfrombyte == nil{
		return nil
	}
	
	idtobyte := txutil.AddrHelper.DecodeUserid(to)
	if idtobyte == nil{
		return nil
	}	
	
	idtotal := append(idfrombyte, idtobyte...)
	
	shabyte := sha256.Sum256(append(idtotal, nounce...))
	return shabyte[:]
}

//so we get three types return: true and no error indicate we definitely get the exist nounce,
//false and no error indicate we definitely not get the exist nounce or it has been expired,
//false and error indicate we could not know the nounce exist or not and it is on your risk to continue
func (m *NounceManager) CheckfundNounce(stub shim.ChaincodeStubInterface, from string, to string, nounce []byte) (bool, error){
	
	m.nouncekey = GenfundNounce(from, to, nounce)
	if m.nouncekey == nil{
		return false, errors.New("Could not get func nounce key")
	}
	
	nouncekey := GenFuncNounceKeyStr(m.nouncekey)
	logger.Debug("fund tx nounce:", nouncekey)
	m.stub = stub
	
	data, err := stub.GetState(nouncekey)
	if err != nil{
		return false, err
	}
	
	if data != nil{	
		//just check the data ...
		nouncedata := &pb.NounceData{}
		err = proto.Unmarshal(data, nouncedata)
		if err == nil{
			logger.Warning("May encounter a replay tx, original is in", nouncedata.NounceTime)			
		}else{
			logger.Error("Recorded nounce is invalid:", err)
		}
		
		return true, err
	}
	
	return false, nil
	
}

func (m *NounceManager) SavefundNounce(stub shim.ChaincodeStubInterface, from *pb.UserData, to *pb.UserData){
	if len(m.nouncekey) == 0{
		return
	}
	
	nouncedata := &pb.NounceData{Txid: stub.GetTxID(), NounceTime: m.Tsnow,
		FromLast: from.LastFund, ToLast: to.LastFund}
	data, err := proto.Marshal(nouncedata)
	if err == nil{
		m.stub.PutState(GenFuncNounceKeyStr(m.nouncekey), data)
	}else{
		logger.Error("Marshal nounce data fail!", err)
	}
	
}

