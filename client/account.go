package client

import (
	"fmt"
	"errors"
	"strings"
	
	"gamecenter.mobi/paicode/wallet"
	txutil "gamecenter.mobi/paicode/transactions"
)

type accountManager struct{
	KeyMgr wallet.KeyManager
}


//generate privatekey: <remark>
func (m* accountManager) GenPrivkey(args ...string) (string, error){
	if len(args) > 1{
		return "", errors.New(fmt.Sprint("Could not recognize", args[1:]))
	}
	
	var remark string		
	if len(args) == 0{
		remark = RandStringRunes(16)
	}else{
		remark = args[0]
	}
	
	k, err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		return remark, err
	}
	
	m.KeyMgr.AddPrivKey(remark, k)
	
	return remark, nil
		
}

//dump privatekey from [remark]
func (m* accountManager) DumpPrivkey(args ...string) (string, error){
	if len(args) != 1{
		return "", errors.New("Invalid remark")
	}
	
	k, err := m.KeyMgr.LoadPrivKey(args[0])
	if err != nil{
		return "", err
	}
	
	return k.DumpPrivkey()
}

//dump privatekey from [remark]
func (m* accountManager) DelPrivkey(args ...string)  error{
	if len(args) != 1{
		return errors.New("Invalid arguments")
	}
	
	return m.KeyMgr.RemovePrivKey(args[0])
}

//get address from [remark]
func (m* accountManager) GetAddress(args ...string) (string, error){
	if len(args) != 1{
		return "", errors.New("Invalid remark")
	}
	
	k, err := m.KeyMgr.LoadPrivKey(args[0])
	if err != nil{
		return "", err
	}

	addr := txutil.AddrHelper.GetUserId(&k.K.PublicKey)
	if len(addr) == 0{
		return addr, errors.New("Can't generate userid")
	}
	
	return addr, nil
}

//list all keys in manager with remark and address
func (m* accountManager) ListKeyData(args ...string) [][2]string{
	if len(args) != 0{
		return nil
	}
	
	kmap, err := m.KeyMgr.ListAll()
	if err != nil{
		return nil
	}
	
	ret := make([][2]string, len(kmap))
	var i int = 0
	for k, v := range kmap{
		ret[i] = [2]string{k, txutil.AddrHelper.GetUserId(&v.K.PublicKey)}
		i++
	}
	
	return ret
}

//[import string], <remark>
func (m* accountManager) ImportPrivkey(args ...string) (string, error){
	
	if len(args) == 0{
		return "", errors.New("Need import string")
	}	
	
	if len(args) > 2{
		return "", errors.New(fmt.Sprint("Could not recognize", args[2:]))
	}
	
	var remark string		
	if len(args) == 1{
		remark = RandStringRunes(16)
	}else{
		remark = args[1]
	}	
	
	k, err := wallet.DefaultWallet.ImportPrivKey(args[0])
	if err != nil{
		return "", err
	}
	
	//check duplication
	addr := txutil.AddrHelper.GetUserId(&k.K.PublicKey)
	kmap, err := m.KeyMgr.ListAll()
	if err != nil{
		return "", err
	}
	for _, v := range kmap{
		if strings.Compare(addr, txutil.AddrHelper.GetUserId(&v.K.PublicKey)) == 0{
			return "", errors.New("Duplicate key")
		}
	} 
	
	m.KeyMgr.AddPrivKey(remark, k)
	
	return remark, nil
}

