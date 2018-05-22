package client

import (
    "testing"
    "strings"
    "gamecenter.mobi/paicode/wallet"
)

func TestMgrPrivk(t *testing.T) {
	mgr := accountManager{wallet.CreateSimpleManager("")}
	mgr2 := accountManager{wallet.CreateSimpleManager("")}
	
	_, err := mgr.GenPrivkey("test1")
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = mgr.GenPrivkey("test11")
	if err != nil{
		t.Fatal(err)
	}	
	
	dstr1, err := mgr.DumpPrivkey("test1")
	if err != nil{
		t.Fatal(err)
	}
	
	dstr11, err := mgr.DumpPrivkey("test11")
	if err != nil{
		t.Fatal(err)
	}	
	
	_, err = mgr.DumpPrivkey("test2")
	if err == nil{
		t.Fatal("Dump unexist key")
	}
	
	_, err = mgr2.ImportPrivkey(dstr1)
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = mgr2.ImportPrivkey(dstr11, "test2")
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = mgr.ImportPrivkey(dstr11)
	if err == nil{
		t.Fatal("Import duplicated key")
	}	
	
	dstr2, err := mgr2.DumpPrivkey("test2")
	if err != nil{
		t.Fatal(err)
	}
	
	if strings.Compare(dstr11, dstr2) != 0{
		t.Fatal("Dumped key not identical")
	}
	
	list := mgr.ListKeyData()
	if len(list) != 2{
		t.Fatal("Not expected count for keys")
	}
	
	t.Log(list)
	
	_, err = mgr2.GetAddress("test2")
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = mgr2.GetAddress("test3")
	if err == nil{
		t.Fatal("Get unexist address")
	}
		
}

