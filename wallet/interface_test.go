package wallet

import (
    "testing"
    "bytes"
    "fmt"
    "crypto/rand"
    "math/big"
)

type sampleSets struct{
	sampleCnt  int	
	existGroup KeyManager
	nonexistGroup KeyManager
}

func (s *sampleSets) prepare() error{
	//simple manager should be robust and can be used for test samples
	s.existGroup = CreateSimpleManager("")
	s.nonexistGroup = CreateSimpleManager("")
	
	if s.sampleCnt == 0{
		s.sampleCnt = 5000
	}
	
	randlimit := big.NewInt(0xFFFFFFFF)
	
	for i := 0; i < s.sampleCnt; i++{
		k, err := DefaultWallet.GeneratePrivKey()
		if err != nil{
			return err
		}
		
		rbn, err := rand.Int(rand.Reader, randlimit)
		if err != nil{
			return err
		}
				
		s.existGroup.AddPrivKey(fmt.Sprintf("TestSet%v%v", i, rbn.Int64()), k)
		s.nonexistGroup.AddPrivKey(fmt.Sprintf("TestSetNot%v%v", i, rbn.Int64()), k)
	}
	
	return nil
}

func testChargeProcess(t *testing.T, target KeyManager, sample *sampleSets) {
	vm, err := sample.existGroup.ListAll()
	if err != nil{
		t.Fatal(err)
	}
	
	for k, v := range vm{
		target.AddPrivKey(k, v)
	}	
}

func testStandardProcess(t *testing.T, target KeyManager, sample *sampleSets) {
	
	vm, err := sample.existGroup.ListAll()
	if err != nil{
		t.Fatal(err)
	}
	
	for k, v := range vm{
		v1, err := target.LoadPrivKey(k)
		if err != nil{
			t.Fatal(err)
		}
		
		if bytes.Compare(v1.K.D.Bytes(), v.K.D.Bytes()) != 0{
			t.Fatal("Value not match for key", k)
		}
	}
	
	nvm, err := sample.nonexistGroup.ListAll()
	if err != nil{
		t.Fatal(err)
	}
	
	for k, _ := range nvm{
		_, err := target.LoadPrivKey(k)
		if err == nil{
			t.Fatal("Load key should not exist")
		}
	}		
}

func TestSimpleManager(t *testing.T) {
	
	test1 := CreateSimpleManager("Test.dat")
	
	set1 := &sampleSets{sampleCnt: 1000}
	err := set1.prepare()
	
	if err != nil{
		t.Fatal(err)		
	}
	
	set2 := &sampleSets{sampleCnt: 2000}
	err = set2.prepare()
	
	if err != nil{
		t.Fatal(err)		
	}
	
	set3 := &sampleSets{sampleCnt: 4000}
	err = set3.prepare()
	
	if err != nil{
		t.Fatal(err)		
	}		
	
	testChargeProcess(t, test1, set1)
	testChargeProcess(t, test1, set2)
	testStandardProcess(t, test1, set1)
	testChargeProcess(t, test1, set3)
	testStandardProcess(t, test1, set3)
	testStandardProcess(t, test1, set2)
	testStandardProcess(t, test1, set1)
	
	err = test1.Persist()
	
	if err != nil{
		t.Fatal(err)		
	}
	
	test2 := CreateSimpleManager("Test.dat")
	err = test2.Load()
	if err != nil{
		t.Fatal(err)		
	}
	
	testStandardProcess(t, test2, set2)
	testStandardProcess(t, test2, set3)
	testStandardProcess(t, test2, set1)		
	
	test3 := CreateSimpleManager("Test.dat")
	testChargeProcess(t, test3, set2)
	err = test3.Load()
	if err != nil{
		t.Fatal(err)		
	}
	
	testStandardProcess(t, test3, set1)
	testStandardProcess(t, test3, set3)
	testStandardProcess(t, test3, set2)		
	
}

