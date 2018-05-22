package main

import (
	"testing"
	_ "strings"
	_ "bytes"
	_ "errors"
	"crypto/ecdsa"
	_ "crypto/elliptic"
	_ "crypto/rand"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	proto "github.com/golang/protobuf/proto"
	
	paicrypto "gamecenter.mobi/paicode/crypto"
	"gamecenter.mobi/paicode/wallet"
	pb "gamecenter.mobi/paicode/protos"
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	txutil "gamecenter.mobi/paicode/transactions"	
)

func init(){
	txutil.AddrHelper = txutil.AddressHelper(13)
}


func makeInit(stub *shim.MockStub, total int64, preassign map[string]int64) error{
	
	initset := &pb.InitChaincode{&pb.DeploySetting{true, int32(txutil.AddrHelper), total, total}, nil}
	
	for k, v := range preassign{
		initset.PreassignedUser = append(initset.PreassignedUser, &pb.PreassignData{k, v})
	}
	
	arg, err := txutil.EncodeChaincodeTx(initset)
	if err != nil {
		return err
	}
	
	_, err = stub.MockInit("1", "init", []string{arg})
	if err != nil {
		return err
	}
	
	return nil
}

func checkGlobalPai(t *testing.T, stub *shim.MockStub, expect int64 ) {
	buf, ok := stub.State[global_setting_entry]
	if !ok{
		t.Fatal("No global status")
	}
	
	ret := &pb.DeploySetting{}
	err := proto.Unmarshal(buf, ret)
	if err != nil{
		t.Fatal("Unmarshal fail", err)
	}
	
	if !ret.DebugMode || int32(txutil.AddrHelper) != ret.NetworkCode || expect != ret.UnassignedPais{
		t.Fatal("Not correct global setting", ret)
	}
}

func checkUser(t *testing.T, stub *shim.MockStub, uid string, expect int64){
	buf, ok := stub.State[uid]
	if !ok{
		t.Fatal("No user", uid)
	}
	
	ret := &pb.UserData{}
	err := proto.Unmarshal(buf, ret)
	if err != nil{
		t.Fatal("Unmarshal fail", err)
	}
	
	if ret.Pais != expect{
		t.Fatal("Not correct pais for user", uid, ret.Pais, expect)
	}
}

func TestPaichaincode_Init(t *testing.T) {
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	err := makeInit(stub, 100000, map[string]int64{})
	if err != nil{
		t.Fatal(err)
	}

	checkGlobalPai(t, stub, 100000)
}

func TestPaichaincode_InitPreassign(t *testing.T) {
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	err := makeInit(stub, 100000, map[string]int64{"dummy1": 50000, "dummy2": 10})
	if err != nil{
		t.Fatal(err)
	}

	checkGlobalPai(t, stub, 49990)
	checkUser(t, stub, "dummy1", 50000)
	checkUser(t, stub, "dummy2", 10)
}

func TestPaichaincode_InitPreset(t *testing.T) {
	
	if int(txutil.AddrHelper) != 13{
		t.Skip("Not match network code for this test, require 13")
	}
	
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	_, err := stub.MockInit("1", "init", []string{"CgwIARANUMCEPVjAhD0SKAoiRFhQbFJtNGxxeTNzTVEyeHVrNzJWUExPYkp3TGlwbTdWQRCgwh4="})
	if err != nil {
		t.Fatal(err)
	}

	checkGlobalPai(t, stub, 500000)
	checkUser(t, stub, "DXPlRm4lqy3sMQ2xuk72VPLObJwLipm7VA", 500000)
}

type privKey struct{
	k 			*ecdsa.PrivateKey
	underlyingK	*paicrypto.ECDSAPriv
} 

func producePrivk(count int) (ret []*wallet.Privkey){
	
	ret = make([]*wallet.Privkey, count)
	for i, _ := range ret{
		k , _ := wallet.DefaultWallet.GeneratePrivKey()
		if k == nil{
			ret = nil
			return
		}
		
		ret[i] = k
	}
	
	return
} 

func confirmUser(t *testing.T, stub *shim.MockStub, privk *wallet.Privkey){
	pd := &txutil.UserTxProducer{PrivKey: privk.K}
	
	args, err := pd.MakeArguments(&pb.RegPublicKey{privk.GenPublicKeyMsg()})
	
	if err != nil{
		t.Fatal(err)		
	}
	
	_, err = stub.MockInvoke("confirmTest", tx.UserRegPublicKey, args)
	if err != nil{
		t.Fatal(err)
	}
}

func testFailConfirmUser(t *testing.T, stub *shim.MockStub, privk *wallet.Privkey, wrong_privk *wallet.Privkey){
	pd := &txutil.UserTxProducer{PrivKey: privk.K}
	
	args, err := pd.MakeArguments(&pb.RegPublicKey{wrong_privk.GenPublicKeyMsg()})
	
	if err != nil{
		t.Fatal(err)		
	}
	
	_, err = stub.MockInvoke("confirmTest", tx.UserRegPublicKey, args)
	if err == nil{
		t.Fatal("Not fail wrong user reg")
	}
	
	t.Log(err)
}

func confirmPais(t *testing.T, stub *shim.MockStub, uid string, expect int64 ){
	
	v, ok := stub.State[uid]
	if !ok{
		t.Fatal("No user id", uid)
	}
	
	ud := &pb.UserData{}
	err := proto.Unmarshal(v, ud)	
	if err != nil{
		t.Fatal(err)		
	}
	
	if ud.Pais != expect{
		t.Fatal("Pai not match", ud.Pais, expect)
	}

}

func TestPaichaincode_FundTx(t *testing.T) {
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	keys := producePrivk(3)
	if len(keys) < 3{
		t.Fatal("Produce keys fail")
	}

	ids := [3]string{
		txutil.AddrHelper.GetUserId(&keys[0].K.PublicKey),
		txutil.AddrHelper.GetUserId(&keys[1].K.PublicKey),
		txutil.AddrHelper.GetUserId(&keys[2].K.PublicKey)}

	err := makeInit(stub, 100000, map[string]int64{
			ids[0]: 50000, 
			ids[1]: 10,
			ids[2]: 99})
	
	if err != nil{
		t.Fatal(err)
	}

	confirmUser(t, stub, keys[0])
	confirmUser(t, stub, keys[2])
	testFailConfirmUser(t, stub,  keys[1],  keys[2])

	nounce1 := []byte{42,42,42}

	tx1 := &tx.FundTx{tx.FundTxData{ids[1], 100}, nounce1, false, 0}
	args, err := tx1.MakeTransaction(keys[0].K) 
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err != nil{
		t.Fatal(err)
	}
	
	confirmPais(t, stub, ids[0], 49900)
	confirmPais(t, stub, ids[1], 110)
	
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err == nil{
		t.Fatal("Do not recognize duplicated")
	}
	t.Log(err)
	
	confirmPais(t, stub, ids[0], 49900)	
	confirmPais(t, stub, ids[1], 110)
	
	tx2 := &tx.FundTx{tx.FundTxData{ids[0], 10}, nounce1, false, 0}
	args, err = tx2.MakeTransaction(keys[1].K) 
	if err != nil{
		t.Fatal(err)
	}	
	
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err == nil{
		t.Fatal("Do not recognize duplicated")
	}	
	if err == nil{
		t.Fatal("Do not recognize unregistered user")
	}
	t.Log(err)
	
	args, err = tx2.MakeTransaction(keys[2].K) 
	if err != nil{
		t.Fatal(err)
	}	
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err != nil{
		t.Fatal(err)
	}	
		
	confirmPais(t, stub, ids[0], 49910)	
	confirmPais(t, stub, ids[2], 89)	
	
	args, err = tx2.MakeTransaction(keys[0].K) 
	if err != nil{
		t.Fatal(err)
	}	
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err == nil{
		t.Fatal("Do not recognize self-funding")
	}		
	t.Log(err)
	
	nounce2 := []byte{42,42,42,42}
	tx3 := &tx.FundTx{tx.FundTxData{ids[1], 10000}, nounce2, false, 0}
	args, err = tx3.MakeTransaction(keys[2].K) 
	if err != nil{
		t.Fatal(err)
	}
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err == nil{
		t.Fatal("Do not recognize unenough pais")
	}		
	t.Log(err)	
	
	args, err = tx3.MakeTransaction(keys[0].K) 
	if err != nil{
		t.Fatal(err)
	}
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err != nil{
		t.Fatal(err)
	}	
	
	confirmPais(t, stub, ids[0], 39910)	
	confirmPais(t, stub, ids[1], 10110)		
	
	tx4 := &tx.FundTx{tx.FundTxData{ids[1], 10000}, nil, false, 0}
	args, err = tx4.MakeTransaction(keys[0].K) 
	if err != nil{
		t.Fatal(err)
	}
	
	argsfake, _ := tx4.MakeTransaction(keys[1].K) 
	argsfake[0] = args[0] //try to breach 0's pai by another signature
	_, err = stub.MockInvoke("testFund", tx.UserFund, argsfake)
	if err == nil{
		t.Fatal("Do not recognize fake signature")
	}		
	t.Log(err)		
	
	_, err = stub.MockInvoke("testFund", tx.UserFund, args)
	if err != nil{
		t.Fatal(err)
	}
	
	confirmPais(t, stub, ids[0], 29910)	
	confirmPais(t, stub, ids[1], 20110)		
}


