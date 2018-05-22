package transactions

import (
	"crypto/ecdsa"
	"math/big"
	
	pb "gamecenter.mobi/paicode/protos"
	paicrypto "gamecenter.mobi/paicode/crypto"

)

type PublicKey  pb.PublicKey

func (pk PublicKey) ECDSAPublicKey() *ecdsa.PublicKey{
	curve, err := paicrypto.GetEC(int(pk.Curvetype))
	if err != nil{
		return nil
	}
	
	return &ecdsa.PublicKey{curve, big.NewInt(0).SetBytes(pk.P.X),
		big.NewInt(0).SetBytes(pk.P.Y)}	
}


func MakePbFromPrivKey(priv *paicrypto.ECDSAPriv) (*pb.PublicKey, error){
	
	pk, err := priv.Apply()
	if err != nil {
		return nil, err
	}
	
	return &pb.PublicKey{int32(priv.CurveType), &pb.ECPoint{pk.X.Bytes(), pk.Y.Bytes()}}, nil
	
}

