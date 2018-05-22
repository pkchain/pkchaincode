package transaction

import (
	"fmt"
	"strings"
	"errors"
	"crypto/ecdsa"
	"time"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	
	txutil "gamecenter.mobi/paicode/transactions"
	pb 	   "gamecenter.mobi/paicode/protos"
	persistpb "gamecenter.mobi/paicode/protos"
)

type FundTxData struct{
	To   string
	Amount uint	
}

type FundTx struct{
	FundTxData
	Nounce []byte
	Invoked bool
	InvokedCode uint
}

func (f *FundTxData) fill(v *pb.Funddata) {
	
	f.Amount = uint(v.Pai)
	f.To = v.ToUserId
}

func (f *FundTx) fill(v interface{}) error{
	switch data := v.(type){
		case *pb.UserTxHeader:
		f.Nounce = data.Nounce
		case *pb.Fund:
		f.Invoked = data.InvokeChaincode != 0
		f.InvokedCode = uint(data.InvokeChaincode)
		switch inndata := data.D.(type){
			case *pb.Fund_Userfund:
			f.FundTxData.fill(inndata.Userfund)
			default://simply omit
		}
		case *pb.Funddata:
		f.FundTxData.fill(data)
		default:
		return errors.New(fmt.Sprint("encounter unexpected type as %T", data))
	}
	
	return nil
}

func (f *FundTx) Parse(pk *ecdsa.PublicKey, args []string) error{
	
	cs := txutil.UserTxConsumer{PublicKey: pk}

	if len(args) == 3{
		v := &pb.Fund{}
		err := cs.ParseArguments(args, v)
		if err != nil{
			return err
		}
		
		err = f.fill(v)
		if err != nil{
			return err
		}
				
		if f.Invoked {
			return errors.New("A invoked fund tx with not enough arguments")
		}
		
		f.FundTxData.fill(v.D.(*pb.Fund_Userfund).Userfund) 
		
	}else{
		v1 := &pb.Fund{}
		v2 := &pb.Funddata{}
		err := cs.ParseArguments(args, v1, v2)
		if err != nil{
			return err
		}
		
		err = f.fill(v1)
		if err != nil{
			return err
		}

		if !f.Invoked {
			return errors.New("Not a invoked fund tx")
		}
		
		err = f.fill(v2)
		if err != nil{
			return err
		}
	}
	
	return f.fill(cs.HeaderCache)	
}

func (f *FundTx) MakeTransaction(privk *ecdsa.PrivateKey) ([]string, error){
	
	pd := txutil.UserTxProducer{PrivKey: privk, Nounce: f.Nounce}
	fmain := &pb.Fund{}
	fdata := &pb.Funddata{uint32(f.Amount), f.To}
	
	if f.Invoked {
		//todo: JUST for debugging 
		fmain.InvokeChaincode = uint32(f.InvokedCode)
		fmain.D = &pb.Fund_Null{}
		return pd.MakeArguments(fmain, fdata)
	}else{
		fmain.InvokeChaincode = 0
		fmain.D = &pb.Fund_Userfund{fdata}
		return pd.MakeArguments(fmain)
	}
}

type fundHandler struct{
	
}

func init(){
	UserTxMap[UserFund] = &fundHandler{}
}

func acquireTsNow(stub shim.ChaincodeStubInterface) *timestamp.Timestamp{
	tsnow, err := stub.GetTxTimestamp()
	if tsnow == nil || err != nil{
		logger.Debug("Can't not get timestamp from stub", err)
		return &timestamp.Timestamp{Seconds: time.Now().Unix(), Nanos: 0}
	}else{
		//notice the tsnowx has different type (vendored) with tsnow
		return &timestamp.Timestamp{Seconds: tsnow.Seconds, Nanos: tsnow.Nanos}
	}	
}

func (f *fundHandler) Handle(uid string, ud *persistpb.UserData, stub shim.ChaincodeStubInterface, 
	args []string) (outud map[string]*persistpb.UserData, err error) {
	
	fdetail := &FundTx{}
	err = fdetail.Parse((*txutil.PublicKey) (ud.Pk).ECDSAPublicKey(), args)
	
	if err != nil{
		return
	}
	
	/*check invoking ...*/
	if fdetail.Invoked {
		var isauth bool = false
		logger.Debug("fund tx is invoked") 
		for uauth := range ud.Authcodes{
			if uint(uauth) == fdetail.InvokedCode{
				isauth = true
				break;
			}
		}
		
		if !isauth{
			err = errors.New(fmt.Sprint("Tx is invoked by an non-auth chaincode", fdetail.InvokedCode))
		}
	}
	
	if strings.Compare(uid, fdetail.To) == 0{
		err = errors.New("Can't make funding to yourself")
		return
	} 
	
	/*then checking validation of pai ...*/
	if ud.Pais < int64(fdetail.Amount){
		err = errors.New("User has not enough pais to pay!")
		return
	}
	
	ncMgr := &NounceManager{Tsnow: acquireTsNow(stub)}
	
	if b, errx := ncMgr.CheckfundNounce(stub, uid, fdetail.To, fdetail.Nounce); b{
		err = errors.New("A fund tx has been recorded recently")
		return
	}else if errx != nil{
		logger.Warning("Could not check fund nounce:", errx, ",we just continue")
	}	
	
	//now get the target id
	var data []byte
	data, err = stub.GetState(fdetail.To)
	if err != nil{
		return
	}
	
	toUd := &pb.UserData{LastFund: &pb.FuncRecord{nil, false}}
	if data != nil{		
		err = proto.Unmarshal(data, toUd)
		if err != nil{
			logger.Error("Data for user", fdetail.To, "can not be parsed!")
			return
		}		
	}else{
		//index must be added, for safety we set the pai number
		toUd.Pais = 0
	}

	//we must save nounce first
	ncMgr.SavefundNounce(stub, ud, toUd)	
	
	//finally, update for transaction
	ud.Pais -= int64(fdetail.Amount)
	toUd.Pais += int64(fdetail.Amount)

	ud.LastActive = ncMgr.Tsnow
	toUd.LastActive = ncMgr.Tsnow
	
	ud.LastFund = &pb.FuncRecord{ncMgr.nouncekey, true}
	toUd.LastFund = &pb.FuncRecord{ncMgr.nouncekey, false}	
	
	//all the user data will be written back
	outud = map[string]*persistpb.UserData{uid: ud, fdetail.To: toUd}
	
	logger.Info("Fund transaction of ", fdetail.Amount, "pais from", uid, "to", fdetail.To)
		
	return
} 
