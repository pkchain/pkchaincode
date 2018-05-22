package transactions

import (
	"errors"
	"hash"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	
	"github.com/golang/protobuf/proto"
	pb "gamecenter.mobi/paicode/protos"

)


//all user's transaction for chaincode should has following styles:
//1. A pb.UserTxHeader message
//2. Any additional messages, will be signed with UserTxHeader together
//3. A pb.Signature message
//.... and any other message which is unsigned

type UserTxProducer struct{
	PrivKey *ecdsa.PrivateKey
	customAddrHelper *AddressHelper
	Nounce  []byte
}

const(
	SecurityNounceLen int = 16
	TxTimeStampExpireSec int = 604800 //7 days
)

func (m *UserTxProducer) helpMakingNounce(){
	m.Nounce = make([]byte, SecurityNounceLen)
	rand.Reader.Read(m.Nounce)
}

//the first msg will be signed with UserTxHeader together and the other will be only encode
func (m *UserTxProducer) MakeArguments(msgs ...proto.Message) (args []string, err error){
	if msgs == nil{
		return nil, errors.New("Require at least one messages")
	}
	
	hasher := sha256.New()	
	args = make([]string, len(msgs) + 2)
	
	var userid string
	if m.customAddrHelper == nil{
		userid = AddrHelper.GetUserId(&m.PrivKey.PublicKey)
	}else{
		userid = m.customAddrHelper.GetUserId(&m.PrivKey.PublicKey)
	}
	
	if m.Nounce == nil{
		m.helpMakingNounce()
	}
	
	args[0], err = EncodeChaincodeTx(&pb.UserTxHeader{FundId: userid, Nounce: m.Nounce}, hasher)
	if err != nil{
		return
	}
	
	args[1], err = EncodeChaincodeTx(msgs[0], hasher)
	if err != nil{
		return
	}
	
	rx, ry, errx := ecdsa.Sign(rand.Reader, m.PrivKey, hasher.Sum(nil))
	if errx != nil{
		err = errx
		return
	}
	args[2], err = EncodeChaincodeTx(&pb.Signature{P: &pb.ECPoint{rx.Bytes(), ry.Bytes()}})
	if err != nil{
		return
	}		

	for i, msg := range msgs[1:]{
		args[i + 3], err = EncodeChaincodeTx(msg)
		if err != nil{
			return
		}		
	}
	
	return
}

type UserTxConsumer struct{
	PublicKey *ecdsa.PublicKey
	HeaderCache *pb.UserTxHeader
	hashCache hash.Hash
}

func (m UserTxConsumer) GetUserId() string{
	if m.HeaderCache == nil{
		return ""
	}
	
	return m.HeaderCache.FundId
}

func (m UserTxConsumer) GetTxNounce() []byte{
	if m.HeaderCache == nil{
		return nil
	}
	
	return m.HeaderCache.Nounce
}

//Consumer often require a 2-phase parsing to obtain its publickey
func (m *UserTxConsumer) Reset() {
	m.HeaderCache = nil
	m.hashCache = nil
}

//Consumer often require a 2-phase parsing to obtain its publickey
func (m *UserTxConsumer) ParseArgumentsFirst(args []string) error{

	if args == nil{
		return errors.New("Wrong empty arguments")
	}
	
	m.hashCache = sha256.New()
	m.HeaderCache = new(pb.UserTxHeader)
	return DecodeChaincodeTx(args[0], m.HeaderCache, m.hashCache)
}

//args will be consumed until all msgs is decoded, including the argument of signature 
func (m *UserTxConsumer) ParseArguments(args []string, msgs ...proto.Message) error{
	
	if msgs == nil{
		return errors.New("Require at least one messages")
	}	
	
	if len(args) + 2 < len(msgs){
		return errors.New("No enough arguments for given messages")
	}
	
	var err error
	if m.hashCache == nil || m.HeaderCache == nil{
		err = m.ParseArgumentsFirst(args)
		if err != nil{
			return err
		}		
	} 
	
	hasher := m.hashCache
	
	err = DecodeChaincodeTx(args[1], msgs[0], hasher)
	if err != nil{
		return err
	}	
	
	if m.PublicKey != nil{
		
		sign := &pb.Signature{}
		err = DecodeChaincodeTx(args[2], sign)
		if err != nil{
			return err
		}		
		
		x, y := big.NewInt(0).SetBytes(sign.P.X), big.NewInt(0).SetBytes(sign.P.Y)
		if !ecdsa.Verify(m.PublicKey, hasher.Sum(nil), x, y){
			return errors.New("Signature not match")
		}
	}
	
	for i, msg := range msgs[1:]{
		err = DecodeChaincodeTx(args[i + 3], msg)
		if err != nil{
			return err
		}		
	}		
	
	return nil
}



