package client

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	
	"gamecenter.mobi/paicode/wallet"
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	txutil "gamecenter.mobi/paicode/transactions"
	pb "gamecenter.mobi/paicode/protos"
	
	"github.com/hyperledger/fabric/peerex"
)

type rpcManager struct{
	PrivKey    *wallet.Privkey
	Rpcbuilder *peerex.RpcBuilder
}

const(
	fundNounceMaxLen int = 256
)

func MakeInitParam(debugmode bool, args ...string) (string, error){
	
	if len(args) < 2{
		return "", errors.New("No enough arguments")
	}
	
	total, err := strconv.Atoi(args[0])
	if err != nil{
		return "", err
	}	
	
	initset := &pb.InitChaincode{&pb.DeploySetting{debugmode, int32(txutil.AddrHelper), int64(total), int64(total)}, nil}
	
	for _, assignstr := range args[1:]{
		ret := strings.Split(assignstr, ":")
		if len(ret) != 2{
			return "", fmt.Errorf("Invalid assign str %s", assignstr)
		}
		
		if amount, err := strconv.Atoi(ret[1]); err != nil{
			return "", fmt.Errorf("Invalid assign amount: %s (in %s)", err.Error(), assignstr)
		}else{
			initset.PreassignedUser = append(initset.PreassignedUser, &pb.PreassignData{ret[0], int64(amount)})
		}
		
	}
	
	return txutil.EncodeChaincodeTx(initset)
}

//funding: <to:addr> <amount> [message]
func (m* rpcManager) Fund(args ...string) (string, error){
	if len(args) < 2{
		return "", errors.New("No required arguments")
	}
		
	b, err := txutil.AddrHelper.VerifyUserId(args[0])
	if !b{
		return "", err
	}
	
	i, err := strconv.Atoi(args[1])
	if err != nil{
		return "", err
	}
	
	if i < 0{
		return "", errors.New("Invalid amount")
	}
	
	fund := &tx.FundTx{FundTxData: tx.FundTxData{args[0], uint(i)}, Invoked: false}
	
	if len(args) == 3{
		if len(args[2]) > fundNounceMaxLen{
			return "", errors.New(fmt.Sprint("message is too long, should not exceed", fundNounceMaxLen, "chars"))
		}
		
		fund.Nounce = []byte(args[2])
	}
	
	rpcargs, err := fund.MakeTransaction(m.PrivKey.K)
	if err != nil{
		return "", err
	}	
	
	m.Rpcbuilder.Function = tx.UserFund
	return m.Rpcbuilder.Fire(rpcargs)
		
}

//registry: <no input>
func (m* rpcManager) Registry(args ...string) (string, error){
	if len(args) != 0{
		return "", errors.New("Not require arguments")
	}
	
	if m.PrivKey == nil || !m.PrivKey.IsValid() {
		return "", errors.New("Key is not applied")		
	}
	
	pd := &txutil.UserTxProducer{PrivKey: m.PrivKey.K}
	regmsg := &pb.RegPublicKey{m.PrivKey.GenPublicKeyMsg()}
	
	rpcargs, err := pd.MakeArguments(regmsg)
	if err != nil{
		return "", err
	}
	
	m.Rpcbuilder.Function = tx.UserRegPublicKey
	return m.Rpcbuilder.Fire(rpcargs)
		
}
