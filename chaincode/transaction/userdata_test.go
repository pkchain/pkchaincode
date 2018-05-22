package transaction

import (
	"testing"
	"bytes"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gamecenter.mobi/paicode/wallet"
	pb "gamecenter.mobi/paicode/protos"
	txutil "gamecenter.mobi/paicode/transactions"	
)

func TestUserDataTx(t *testing.T){

	stub := shim.NewMockStub("DummyTest", nil)	
	
	privk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	
	uid := txutil.AddrHelper.GetUserId(&privk.K.PublicKey)
	
	pd := &txutil.UserTxProducer{PrivKey: privk.K}
	inpk := privk.GenPublicKeyMsg()
	args, err := pd.MakeArguments(&pb.RegPublicKey{inpk})
	
	if err != nil{
		t.Fatal(err)		
	}	
	
	h := &regPublicKeyHandler{}
	dummyRec := &pb.FuncRecord{nil, false}
	
	stub.MockTransactionStart("1")
	out, err := h.Handle(uid, &pb.UserData{0, nil, nil, "Hello region", dummyRec, nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if len(out) != 1 {
		t.Fatal("Invalid output")
	}
	
	if v, ok := out[uid]; !ok || v.Pk == nil || v.Pk.P == nil || bytes.Compare(v.Pk.P.X, privk.K.PublicKey.X.Bytes()) != 0{
		t.Fatal("Error output")
	}
	
	if out[uid].LastActive == nil{
		t.Fatal("Wrong timestamp")
	}
	
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{0, inpk, out[uid].LastActive, "New region", dummyRec, nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}	
	stub.MockTransactionEnd("1")
	
	if len(out) != 1 {
		t.Fatal("Invalid output")
	}	
	
	yaprivk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	
	stub.MockTransactionStart("1")
	out, err = h.Handle(txutil.AddrHelper.GetUserId(&yaprivk.K.PublicKey), &pb.UserData{0, nil, nil, "", dummyRec, nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize error")
	}
	stub.MockTransactionEnd("1")
	
	t.Log(err)
}

