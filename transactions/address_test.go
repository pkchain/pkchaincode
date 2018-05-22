package transactions

import (
	"testing"
	"strings"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
)

func testIdA(t *testing.T, pk1 *ecdsa.PublicKey, pk2 *ecdsa.PublicKey){

	id1 := AddrHelper.GetUserId(pk1)
	id2 := AddrHelper.GetUserId(pk2)
	
	if len(id1) == 0 || len(id2) == 0{
		t.Fatal("Invalid id", id1, "--", id2)
	}
	
	if strings.Compare(id1, id2) != 0{
		t.Fatal("Not identify id", id1, "--", id2)
	}
	
	t.Log(id1)
	
}

func testIdB(t *testing.T, pk1 *ecdsa.PublicKey, pk2 *ecdsa.PublicKey){

	id1 := AddrHelper.GetUserId(pk1)
	id2 := AddrHelper.GetUserId(pk2)
	
	if len(id1) == 0 || len(id2) == 0{
		t.Fatal("Invalid id", id1, "--", id2)
	}
	
	if strings.Compare(id1, id2) == 0{
		t.Fatal("Not unique id", id1, "--", id2)
	}
	
	t.Log(id1, "--", id2)
	
}

func TestDump_Userid(t *testing.T){
	
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	prv1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Make ecdsa key fail", err)
	}	
	prv2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Make ecdsa key fail", err)
	}
	
	testIdA(t, &prv1.PublicKey, &prv1.PublicKey)
	testIdB(t, &prv1.PublicKey, &prv2.PublicKey)
	
	uid := AddrHelper.GetUserId(&prv1.PublicKey)
	decbuf := AddrHelper.DecodeUserid(uid)
	if decbuf == nil{
		t.Fatal("Can't decode address")
	}
	
	if strings.Compare(AddrHelper.GenUserId(decbuf[1:]), uid) != 0{
		t.Fatal("Decoded byte not correct")
	}
	
	for i := 0; i < 5000; i++ {
		prv1 = prv2
		
		testIdA(t, &prv1.PublicKey, &prv2.PublicKey)
		
		prv2, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		
		if err != nil{
			t.Skip("Make ecdsa key fail", err)
		}		
		
		testIdB(t, &prv1.PublicKey, &prv2.PublicKey)
	}
	
	
}

func TestVerify_Userid(t *testing.T){
	
	prv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		t.Skip("Make ecdsa key fail")
	}	
	
	uid := AddrHelper.GetUserId(&prv.PublicKey)
	
	if b, err := AddrHelper.VerifyUserId(uid); !b{
		t.Fatal("verify fail", err)
	}
	
	var yahelper AddressHelper = 0
	
	if b, err := yahelper.VerifyUserId(uid); b{
		t.Fatal("verify error")		
	}else{
		t.Log("verfiy ret", err)
	}
}
