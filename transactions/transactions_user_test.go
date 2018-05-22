package transactions

import (
	"testing"
	"crypto/ecdsa"
	"crypto/rand"	
	"crypto/elliptic"
	
	"github.com/golang/protobuf/proto"
	pb "gamecenter.mobi/paicode/protos"
)

func TestUserTxCoding(t *testing.T){

	fields := [][2]proto.Message{
		{&pb.Fund{InvokeChaincode: 10, D: &pb.Fund_Null{false}}, &pb.Fund{}},
		{&pb.Fund{InvokeChaincode: 0, D: &pb.Fund_Userfund{&pb.Funddata{Pai:1000, ToUserId: "test2"}}}, &pb.Fund{}}}

	privk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}	
	
	pd := UserTxProducer{PrivKey: privk}
	cm := UserTxConsumer{}
	
	args, err := pd.MakeArguments(fields[0][0], fields[1][0])
	if err != nil{
		t.Fatal(err)
	}	
	
	if len(args) != 4{
		t.Fatal("Wrong argument count: not 4")
	}
	
	t.Log("Make arguments", args)
	
	err = cm.ParseArguments(args, fields[0][1], fields[1][1])
	if err != nil{
		t.Fatal(err)
	}
	
	for _, finout := range fields{		
		
		if !proto.Equal(finout[0], finout[1]){
			t.Fatal("Not equal message:", finout[0], finout[1])
		}		
	}	
	
	cm = UserTxConsumer{PublicKey: &privk.PublicKey}
	
	args, err = pd.MakeArguments(fields[0][0], fields[1][0])
	if err != nil{
		t.Fatal(err)
	}	
	
	if len(args) != 4{
		t.Fatal("Wrong argument count: not 4")
	}
	
	t.Log("Make arguments", args)
	
	err = cm.ParseArguments(args, fields[0][1], fields[1][1])
	if err != nil{
		t.Fatal(err)
	}
	
	for _, finout := range fields{		
		
		if !proto.Equal(finout[0], finout[1]){
			t.Fatal("Not equal message:", finout[0], finout[1])
		}		
	}	
	
	args, err = pd.MakeArguments(fields[0][0])
	if err != nil{
		t.Fatal(err)
	}	
	
	if len(args) != 3{
		t.Fatal("Wrong argument count: not 3")
	}
	
	t.Log("Make arguments", args)
	
	cm.Reset()	
	err = cm.ParseArguments(args, fields[0][1])
	if err != nil{
		t.Fatal(err)
	}
	
	for _, finout := range fields[:2]{		
		
		if !proto.Equal(finout[0], finout[1]){
			t.Fatal("Not equal message:", finout[0], finout[1])
		}		
	}	
	
	
}