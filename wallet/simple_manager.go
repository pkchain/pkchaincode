package wallet

import (
	"errors"
	"bytes"
	"io"
	"io/ioutil"
	"sync"
	"encoding/gob"
	paicrypto "gamecenter.mobi/paicode/crypto"
)

type simpleManager struct{
	PersistFile string
	keyData		map[string]*Privkey
	lock		sync.RWMutex
}

func CreateSimpleManager(fpath string) *simpleManager{
	return &simpleManager{PersistFile: fpath, keyData: map[string]*Privkey{}}
}

func (m *simpleManager) AddPrivKey(remark string, privk *Privkey){
	m.lock.Lock()
	defer m.lock.Unlock()
	
	m.keyData[remark] = privk
}

func (m *simpleManager) LoadPrivKey(remark string) (*Privkey, error){
	
	m.lock.RLock()
	defer m.lock.RUnlock()	
	
	k, ok := m.keyData[remark]
	if !ok {
		return nil, errors.New("No this key")
	}
	
	return k, nil
}

func (m *simpleManager) RemovePrivKey(remark string) error{

	_, ok := m.keyData[remark]
	if !ok {
		return errors.New("No this key")
	}	
	
	delete(m.keyData, remark)
	return nil
}

func (m *simpleManager) ListAll() (map[string]*Privkey, error){
	
	m.lock.RLock()
	defer m.lock.RUnlock()		
	
	//we do a deep copy
	copiedmap := map[string]*Privkey{}
	
	for k, v := range m.keyData{
		copiedmap[k] = v
	}
	
	return copiedmap, nil
}

type persistElem struct{
	Key  string
	Dump string
}

const defaultFileName string = "simplewallet.dat"

func (m *simpleManager) Load() (err error){

	if m.keyData == nil{
		m.keyData = map[string]*Privkey{}	
	}
	
	origSize := len(m.keyData)

	var data []byte
	if len(m.PersistFile) == 0{
		data, err = ioutil.ReadFile(defaultFileName)
	}else{
		data, err = ioutil.ReadFile(m.PersistFile)
	}
	
	if err != nil{
		return
	}
	
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	
	v := &persistElem{}
	err = dec.Decode(v)
	for ;err == nil; err = dec.Decode(v) {
		uk, errx := paicrypto.PrivKeyfromString(v.Dump)
		if errx != nil{
			logger.Warning("Restore privkey <", v.Key, "> fail:", errx)
			continue			
		}
		privk, errx := uk.Apply()
		if errx != nil{
			logger.Warning("Get ecdsa privkey <", v.Key, "> fail:", errx)
			continue			
		}
				
		m.keyData[v.Key] = &Privkey{privk, uk}
	}
	
	if err == io.EOF{
		err = nil
	}
	
	logger.Info("Restore",len(m.keyData) - origSize,"keys")
	return
	
}

func (m *simpleManager) Persist() error{
		
	m.lock.RLock()
	defer m.lock.RUnlock()		
		
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	enc := gob.NewEncoder(buf)
	
	var saveSize int = 0
	for k, v := range m.keyData{
		str, err := v.underlyingKey.DumpPrivKey()
		if err != nil{
			logger.Warning("Dump privkey <", k, "> fail:", err)
			continue
		}
				
		err = enc.Encode(&persistElem{k, str})
		
		if err != nil{
			logger.Warning("Encode privkey fail", err)
			continue
		}	
		saveSize++	
	}
	
	logger.Info("Save",saveSize,"keys")
	
	if len(m.PersistFile) == 0{
		return ioutil.WriteFile(defaultFileName, buf.Bytes(), 0666)
	}else{
		return ioutil.WriteFile(m.PersistFile, buf.Bytes(), 0666)
	}	
		
	
}