package crypto

import (
	"fmt"
	"errors"
	"encoding/asn1"
	"encoding/base64"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"	
)

type ECDSAPriv struct {
	CurveType int "tag:10"
	D *big.Int
}

func dumpPrivKey(priv ECDSAPriv) ([]byte, error){

	rbyte, err := asn1.Marshal(priv)
	if err != nil {
		return nil, err
	}
	
	return rbyte, nil
}

func importPrivKey(kb []byte)(*ECDSAPriv, error){
	
	var priv = ECDSAPriv{}
	_, err := asn1.Unmarshal(kb, &priv)
	
	if err != nil {
		return nil, err
	}
	
	return &priv, nil		
}

const(
	
	ECP256_FIPS186 = 1
	ECP256_SEC2k1 = 16
)

func GetEC(curveType int) (elliptic.Curve, error){
	switch(curveType){
		case 1:
			return elliptic.P256(), nil
		default:
			return nil, errors.New(fmt.Sprintf("%d is not a valid curve defination", curveType))		
	}	
}

func PrivKeyfromString(kstr string) (*ECDSAPriv, error){
	
	data, err := base64.StdEncoding.DecodeString(kstr)
	if err != nil {
		return nil, err
	}
	
	return importPrivKey(data)	
}

func (priv ECDSAPriv) DumpPrivKey() (string, error){
	
	rb, err := dumpPrivKey(priv)
	if err != nil{
		return "ERROR", err
	}
	
	return base64.StdEncoding.EncodeToString(rb), nil
}

func (k ECDSAPriv) Apply() (*ecdsa.PrivateKey, error){
	
	if k.D == nil{
		return nil, errors.New(fmt.Sprintf("empty seed field"))	
	}
	
	curve, err := GetEC(int(k.CurveType))
	if err != nil{
		return nil, err
	}
	
	retx, rety := curve.ScalarBaseMult(k.D.Bytes())
	
	return &ecdsa.PrivateKey{  ecdsa.PublicKey{curve, retx, rety}, k.D}, nil
}



