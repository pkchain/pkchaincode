package transactions

import (
	"errors"
	"bytes"
	"math/big"
	"encoding/asn1"
	"encoding/base64"
	"crypto/sha256"
	"crypto/ecdsa"
	"golang.org/x/crypto/ripemd160"
)

type AddressHelper int

var AddrHelper AddressHelper = 1

func getUserHash(pk *ecdsa.PublicKey) ([]byte, error) {
	
	type StandardPk struct{
		X *big.Int
		Y *big.Int
	}
	
	stdpk := StandardPk{X: pk.X, Y: pk.Y}
	
	rawbytes, err := asn1.Marshal(stdpk)
	
	if err != nil{
		return nil, err
	}
	
	hash1pass := sha256.Sum256(rawbytes)
	if len(hash1pass) != sha256.Size{
		return nil, errors.New("Wrong sha256 hashing")
	}
	
	rmd160h := ripemd160.New();
	if nn, err := rmd160h.Write(hash1pass[:]); nn != len(hash1pass) || err != nil{
		return nil, errors.New("Wrong ripemd write")
	}
	
	hash2pass := rmd160h.Sum([]byte{})
	
	if len(hash2pass) != ripemd160.Size{
		return nil, errors.New("Wrong ripemd160 hashing")
	}
	
	return hash2pass, nil
}

const (
	AddressFullByteSize = 25
	AddressPartByteSize = 21
	AddressVerifyCodeSize = 4
)

func (_ AddressHelper) DecodeUserid(id string) []byte{
	
	data, err := base64.RawURLEncoding.DecodeString(id)

	if err != nil {
		return nil
	}
	
	return data[:AddressPartByteSize]	
}

func (_ AddressHelper) GetCheckSum(rb []byte) ([AddressVerifyCodeSize]byte, error){
	
	var ret [AddressVerifyCodeSize]byte
	
	hash1pass := sha256.Sum256(rb)
	if len(hash1pass) != sha256.Size{
		return ret, errors.New("Wrong sha256 hashing 1pass")
	}	
	
	hash2pass := sha256.Sum256(hash1pass[:])
	if len(hash2pass) != sha256.Size{
		return ret, errors.New("Wrong sha256 hashing 2pass")
	}
	
	return [AddressVerifyCodeSize]byte{hash2pass[0], hash2pass[1], hash2pass[2], hash2pass[3]}, nil
}

func (prefix AddressHelper) VerifyUserId(id string) (bool, error){
	data, err := base64.RawURLEncoding.DecodeString(id)

	if err != nil {
		return false, errors.New("Wrong base64 decoding")
	}
	
	if uint8(prefix % 256) != uint8(data[0]){
		return false, errors.New("Different prefix")
	} 
	
	if len(data) != AddressFullByteSize{
		return false, errors.New("Wrong bytes size")
	}
	
	ck, err := prefix.GetCheckSum(data[:AddressPartByteSize])
	if err != nil{
		return false, errors.New("Get checksum fail")
	}
	
	return bytes.Equal(ck[:], data[AddressPartByteSize:]), errors.New("checksum not equal")
}

func (prefix AddressHelper) GenUserId(rb []byte) string{

	fullbytes := make([]byte, 1, AddressFullByteSize)
	fullbytes[0] = uint8(prefix % 256)	
	fullbytes = append(fullbytes, rb...)
	
	if len(fullbytes) != AddressPartByteSize {
		return ""
	}
	
	ck, err := prefix.GetCheckSum(fullbytes)	
	if err != nil{
		return ""
	}
	
	return base64.RawURLEncoding.EncodeToString(append(fullbytes, ck[:]...))
	
}

func (prefix AddressHelper) GetUserId(pk *ecdsa.PublicKey) string{
		
	hashbytes, err := getUserHash(pk)
	if err != nil{
		return ""
	}
	
	return prefix.GenUserId(hashbytes)
}


	
	



