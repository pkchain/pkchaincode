package transactions

import (
	"errors"
	"fmt"
	"hash"
	"encoding/base64"
	
	"github.com/golang/protobuf/proto"

)


func DecodeChaincodeTx(arg string, out proto.Message, hasher... hash.Hash) error{

	data, err := base64.StdEncoding.DecodeString(arg)
	
	if err != nil{
		return errors.New(fmt.Sprint("base64 decode fail:", err))
	}
	
	for _, h := range hasher{
		h.Write(data)
	}
	
	return proto.Unmarshal(data, out)
	
}

func EncodeChaincodeTx(in proto.Message, hasher... hash.Hash) (string, error){

	rb, err := proto.Marshal(in)
	if err != nil{
		return "", err
	}
	
	for _, h := range hasher{
		h.Write(rb)
	}
	
	return base64.StdEncoding.EncodeToString(rb), nil
	
}