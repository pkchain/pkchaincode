package client

import (
	"errors"
	"strings"
	tx "gamecenter.mobi/paicode/chaincode/transaction"
)

//queryuser: <user id>
func (m* rpcManager) QueryUser(args ...string) ([]byte, error){
	if len(args) != 1{
		return nil, errors.New("Require user id")
	}
	
	m.Rpcbuilder.Function = tx.QueryUser
	return m.Rpcbuilder.Query([]string{args[0]})
		
}

//queryrecord: <record key> | <from addr:to addr> <nounce>
func (m* rpcManager) QueryRecord(args ...string) ([]byte, error){
	if len(args) == 0 || len(args) > 2 {
		return nil, errors.New("Not require arguments")
	}
	
	var key string
	if len(args) == 2{
		
		outstr := strings.Split(args[0], ":")
		if len(outstr) != 2{
			return nil, errors.New("Invalid address pair")
		}
		
		kbytes := tx.GenfundNounce(outstr[0], outstr[1], []byte(args[1]))
		if kbytes == nil{
			return nil, errors.New("Invalid address")
		}
		
		key = tx.GenFuncNounceKeyStr(kbytes)
	}else{
		key = args[0]
	}
	
	m.Rpcbuilder.Function = tx.QueryRec
	return m.Rpcbuilder.Query([]string{key})
		
}

func (m* rpcManager) QueryNode(args ...string) ([]byte, error){
	if len(args) != 0{
		return nil, errors.New("Not require arguments")
	}
	
	m.Rpcbuilder.Function = tx.QueryNode
	return m.Rpcbuilder.Query(nil)
		
}

//queryglobal: <no input>
func (m* rpcManager) QueryGlobal(args ...string) ([]byte, error){
	if len(args) != 0{
		return nil, errors.New("Not require arguments")
	}
	
	m.Rpcbuilder.Function = tx.QueryGlobal
	return m.Rpcbuilder.Query(nil)
		
}
