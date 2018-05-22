package transactions

import (
	"testing"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	
	"github.com/golang/protobuf/proto"
	pb "gamecenter.mobi/paicode/protos"
	paicrypto "gamecenter.mobi/paicode/crypto"
)


func TestPublicKey_DumpPb(t *testing.T){

	prvkeystd, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}
	
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	prvk := &paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, prvkeystd.D}
	
	pbdump, err := MakePbFromPrivKey(prvk)
	
	if err != nil {
		t.Fatal(err)
	}
	
	msgbyte, err := proto.Marshal(pbdump)
	if err != nil{
		t.Fatal("Marshal protobuf fail", err)
	}
	
	pbrcv := new(pb.PublicKey)
	
	err = proto.Unmarshal(msgbyte, pbrcv)
	if err != nil{
		t.Fatal("Unmarshal protobuf fail", err)
	}
	
	pk := (*PublicKey) (pbrcv)
	publick2 := pk.ECDSAPublicKey()
	
	if publick2 == nil{
		t.Fatal("Generate public key from NewPublicKey fail")
	}
	
	sx, sy, err := ecdsa.Sign(rand.Reader, prvkeystd, rb)
	if err != nil{
		t.Fatal(err)
	}
	
	if !ecdsa.Verify(publick2, rb, sx, sy) {
		t.Fatal("verify signature with dump publickey fail")
	}	
	
}
