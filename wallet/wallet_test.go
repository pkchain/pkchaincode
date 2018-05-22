package wallet

import (
	"testing"
)

func TestGenkey(t *testing.T){
	
	pk1, err := DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)	
	}

	if pk1 == nil || pk1.K == nil{
		t.Fatal("Empty key")
	}
	

	pk2, err := DefaultWallet.GeneratePrivKey()
	
	if err != nil{
		t.Fatal(err)	
	}

	if pk2 == nil || pk1.K == nil{
		t.Fatal("Empty key")
	}	
	
	if pk1.GenPublicKeyMsg() == nil{
		t.Fatal("Could not gen message")
	}
	
	if pk1.K.D.Cmp(pk2.K.D) == 0{
		t.Fatal("not distributed private key")
	}
}

func TestDumpkey(t *testing.T){
	
	pk1, err := DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)	
	}
	
	str, err := pk1.DumpPrivkey()
	if err != nil{
		t.Fatal(err)	
	}
	
	pk2, err := DefaultWallet.ImportPrivKey(str)
	
	if err != nil{
		t.Fatal(err)	
	}
	
	if pk1.K.D.Cmp(pk2.K.D) != 0{
		t.Fatal("Different key")
	}				
}