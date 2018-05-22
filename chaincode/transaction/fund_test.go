package transaction

import (
	"testing"
	"strings"
	"bytes"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	proto "github.com/golang/protobuf/proto"
	"gamecenter.mobi/paicode/wallet"
	pb "gamecenter.mobi/paicode/protos"
	txutil "gamecenter.mobi/paicode/transactions"		
	
)

func compareTest(tx1 *FundTx, tx2 *FundTx) (bool, error){
	if strings.Compare(tx1.To, tx2.To) != 0 ||
		tx1.Amount != tx2.Amount{
			return false, errors.New("funddata not match")
		}
	
	if tx1.Invoked != tx2.Invoked ||
		(tx1.Invoked && tx1.InvokedCode != tx2.InvokedCode){
			return false, errors.New("Invoke fiels not match")
		}
	
	return true, nil
}

func TestTx_UserFund(t *testing.T){

	privk, err := wallet.DefaultWallet.GeneratePrivKey()
	
	if err != nil{
		t.Fatal(err)
	}	
	
	tx1 := &FundTx{FundTxData{"testB", 100}, nil, false, 0}
	
	args, err := tx1.MakeTransaction(privk.K)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 3 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	txIn := new(FundTx)
	err = txIn.Parse(&privk.K.PublicKey, args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("Output Nounce", txIn.Nounce)
	
	if b, err := compareTest(tx1, txIn); !b{
		t.Fatal(err)
	}
	
	tx2 := &FundTx{FundTxData{"testC", 140}, []byte{44, 44, 44, 44, 44}, true, 13}
	
	args, err = tx2.MakeTransaction(privk.K)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 4 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	err = txIn.Parse(&privk.K.PublicKey, args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("Output Nounce", txIn.Nounce)
	
	if b, err := compareTest(tx2, txIn); !b{
		t.Fatal(err)
	}
	
	if bytes.Compare(tx2.Nounce, txIn.Nounce) != 0{
		t.Fatal("Nounce not match")
	}	
	
}

var stub *shim.MockStub

func TestFundTx(t *testing.T){

	stub = shim.NewMockStub("DummyTest", nil)	
	
	privk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	
	uid := txutil.AddrHelper.GetUserId(&privk.K.PublicKey)
	inpk := privk.GenPublicKeyMsg()
	
	yaprivk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	yauid := txutil.AddrHelper.GetUserId(&yaprivk.K.PublicKey)
	
	tx1 := &FundTx{FundTxData{yauid, 100}, []byte{42,42,42}, false, 0}	
	args, err := tx1.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	h := &fundHandler{}
	
	dummyRec := &pb.FuncRecord{nil, false}
	//first tx
	stub.MockTransactionStart("1")
	out, err := h.Handle(uid, &pb.UserData{1000, inpk, nil, "Heaven", dummyRec, nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if len(out) != 2{
		t.Fatal("Invalid output")
	}
	
	if v, ok := out[uid]; !ok || v.Pais != 900{
		t.Fatal("Invalid funder")
	} 
	
	if v, ok := out[yauid]; !ok || v.Pais != 100{
		t.Fatal("Invalid accepter")
	}
	
	//duplicated tx, should fail
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{900, inpk, nil, "Heaven", dummyRec, nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize duplicated tx")
	}
	t.Log(err)
	stub.MockTransactionEnd("1")
	
	data, err := proto.Marshal(&pb.UserData{50, nil, nil, "Heaven", dummyRec, nil})
	if err != nil{
		t.Fatal(err)
	}
	stub.State[yauid] = data
	
	tx2 := &FundTx{FundTxData{yauid, 100}, []byte{42,42,42,42}, false, 0}	
	args, err = tx2.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}
	
	//another tx, but public key is not match, should fail
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{900, yaprivk.GenPublicKeyMsg(), nil, "Heaven", dummyRec, nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize wrong publickey")
	}
	t.Log(err)
	stub.MockTransactionEnd("1")
	
	//another tx, but no enough pais, should fail
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{10, inpk, nil, "Heaven", dummyRec, nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize not enough pais")
	}
	t.Log(err)
	stub.MockTransactionEnd("1")
	
	//another tx
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{900, inpk, nil, "Heaven", dummyRec, nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if v, ok := out[uid]; !ok || v.Pais != 800{
		t.Fatal("Invalid funder")
	} 
	
	if v, ok := out[yauid]; !ok || v.Pais != 150{
		t.Fatal("Invalid accepter")
	}else{
		//update accepter data
		data, err := proto.Marshal(v)
		if err != nil{
			t.Fatal(err)
		}
		stub.State[yauid] = data		
	}	
	
	tx3 := &FundTx{FundTxData{yauid, 100}, []byte{42,42,42,42,42}, false, 0}	
	args, err = tx3.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	//yet another tx, make pai become zero
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{100, inpk, nil, "Heaven", dummyRec, nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if v, ok := out[uid]; !ok || v.Pais != 0{
		t.Fatal("Invalid funder")
	} 
	
	if v, ok := out[yauid]; !ok || v.Pais != 250{
		t.Fatal("Invalid accepter")
	}	
}

var totalamount int64 = 1000

func verifyOut(t *testing.T, in *pb.UserData, out *pb.UserData, expectin int64) {

	if in == nil || out == nil{
		t.Fatal("Missed output udata")
	}
	
	if in.Pais != expectin || out.Pais != (totalamount - expectin){
		t.Fatal("Wrong output pai")
	}	
	
	if in.LastFund.Nouncekey == nil || out.LastFund.Nouncekey == nil{
		t.Fatal("Missed last record")
	}
	
	if bytes.Compare(in.LastFund.Nouncekey, out.LastFund.Nouncekey) != 0{
		t.Fatal("Wrong last record")
	}

	if !in.LastFund.IsSend || out.LastFund.IsSend {
		t.Fatal("Wrong last record side")
	}
	
	recdata := stub.State[GenFuncNounceKeyStr(in.LastFund.Nouncekey)]
	
	if recdata == nil{
		t.Fatal("No record")
	}
	
	rec := &pb.NounceData{}
	
	err := proto.Unmarshal(recdata, rec)
	if err != nil{
		t.Fatal(err)		
	}	
	
	if strings.Compare(rec.Txid, stub.TxID) != 0{
		t.Fatal("wrong record")
	}
}

func walkthroughRec(t *testing.T, in *pb.UserData){
	
	if in == nil{
		t.Fatal("No data")
	}
	
	k := GenFuncNounceKeyStr(in.LastFund.Nouncekey)
	side := in.LastFund.IsSend
	
	for stub.State[k] != nil{
		
		recdata := stub.State[k]		
		rec := &pb.NounceData{}
	
		err := proto.Unmarshal(recdata, rec)
		if err != nil{
			t.Fatal(err)		
		}		
		
		t.Log("walk through:", rec.Txid)
		
		if side{
			if rec.FromLast.Nouncekey == nil{
				return
			}
			k = GenFuncNounceKeyStr(rec.FromLast.Nouncekey)	
			side = rec.FromLast.IsSend
		}else{
			if rec.FromLast.Nouncekey == nil{
				return
			}
			k = GenFuncNounceKeyStr(rec.ToLast.Nouncekey)	
			side = rec.ToLast.IsSend			
		}
		
	}
	
	t.Fatal("unexpected end")
}

func TestFundTxDetailAndRec(t *testing.T){

	stub = shim.NewMockStub("DummyTest", nil)	
	
	privk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}	
	
	uid := txutil.AddrHelper.GetUserId(&privk.K.PublicKey)
	inpk := privk.GenPublicKeyMsg()
	
	yaprivk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	yauid := txutil.AddrHelper.GetUserId(&yaprivk.K.PublicKey)
	yainpk := yaprivk.GenPublicKeyMsg()
	
	nounce1 := []byte{42,42,42}
	
	tx1 := &FundTx{FundTxData{yauid, 100}, nounce1, false, 0}	
	args, err := tx1.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	h := &fundHandler{}
	
	dummyRec := &pb.FuncRecord{nil, false}
	//first tx
	stub.MockTransactionStart("Txid1")
	out, err := h.Handle(uid, &pb.UserData{totalamount, inpk, nil, "Heaven", dummyRec, nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	
	verifyOut(t, out[uid], out[yauid], 900)
	
	stub.MockTransactionEnd("Txid1")
	
	udata, err := proto.Marshal(out[uid])
	if err != nil{
		t.Fatal(err)
	}	
	stub.State[uid] = udata
	out[yauid].Pk = yainpk
	
	nounce2 := []byte{42,42,42,42}
	
	tx2 := &FundTx{FundTxData{uid, 50}, nounce2, false, 0}	
	args, err = tx2.MakeTransaction(yaprivk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	stub.MockTransactionStart("Txid2")
	out, err = h.Handle(yauid, out[yauid], stub, args)
	if err != nil{
		t.Fatal(err)
	}
	
	verifyOut(t, out[yauid], out[uid], 50)
	
	stub.MockTransactionEnd("Txid2")
	
	udata, err = proto.Marshal(out[yauid])
	if err != nil{
		t.Fatal(err)
	}	
	stub.State[yauid] = udata
	
	nounce3 := []byte{42,42,42,42,42}
	
	tx3 := &FundTx{FundTxData{yauid, 150}, nounce3, false, 0}	
	args, err = tx3.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	stub.MockTransactionStart("Txid3")
	out, err = h.Handle(uid, out[uid], stub, args)
	if err != nil{
		t.Fatal(err)
	}
	
	verifyOut(t, out[uid], out[yauid], 800)
	
	stub.MockTransactionEnd("Txid3")
	
	walkthroughRec(t, out[uid])	
	walkthroughRec(t, out[yauid])	
}

