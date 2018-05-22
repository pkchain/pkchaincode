package transactions

import (
	"testing"
	"bytes"
	"crypto/sha256"
	"github.com/golang/protobuf/proto"
	pb "gamecenter.mobi/paicode/protos"
	
)

func TestEncode_Decode(t *testing.T){
	
	fields := [][2]proto.Message{
		{&pb.UserTxHeader{FundId: "test1", Nounce: []byte{42, 42, 42, 42, 42}}, &pb.UserTxHeader{}},
		{&pb.Fund{InvokeChaincode: 10, D: &pb.Fund_Null{false}}, &pb.Fund{}},
		{&pb.Fund{InvokeChaincode: 0, D: &pb.Fund_Userfund{&pb.Funddata{Pai:1000, ToUserId: "test2"}}}, &pb.Fund{}}}
	
	fin := fields[0][0]
	fout := fields[0][1]
	out, err := EncodeChaincodeTx(fin)
	if err != nil{
		t.Fatal(err)
	}
	
	err = DecodeChaincodeTx(out, fout)
	if err != nil{
		t.Fatal(err)
	}
	
	if !proto.Equal(fin, fout){
		t.Fatal("Not equal message:", fin, fout)
	}	
	
	encH := sha256.New()
	decH := sha256.New()
	
	for _, finout := range fields{
		
		fin := finout[0]
		fout := finout[1]		
		
		out, err := EncodeChaincodeTx(fin, encH)
		if err != nil{
			t.Fatal(err)
		}
		
		err = DecodeChaincodeTx(out, fout, decH)
		if err != nil{
			t.Fatal(err)
		}
		
		if !proto.Equal(fin, fout){
			t.Fatal("Not equal message:", fin, fout)
		}		
	}
	
	if bytes.Compare(encH.Sum(nil), decH.Sum(nil)) != 0{
		t.Fatal("Not equal hashs")
	}
	
	fin = fields[0][0]
	fout = fields[2][1]
	out, err = EncodeChaincodeTx(fin)
	if err != nil{
		t.Fatal(err)
	}
	
	err = DecodeChaincodeTx(out, fout)
	if err == nil{
		t.Fatal("Decode invalid field but not fail")
	}
	
	
}