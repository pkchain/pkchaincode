
package crypto

import (
	"testing"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"	
)

func TestApply_Privkey(t *testing.T){
	
	prvkeystd, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}
	
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	var prvkeyt = ECDSAPriv{ECP256_FIPS186, prvkeystd.D}
	prvkeytapp, err := prvkeyt.Apply()
	
	if err != nil{
		t.Fatal(err)
	}
	
	if prvkeytapp.X.Cmp(prvkeystd.X) != 0 || prvkeytapp.Y.Cmp(prvkeystd.Y) != 0{
		t.Fatal("Unmatch public key:", prvkeytapp.X.Text(16), prvkeystd.X.Text(16), 
			prvkeytapp.Y.Text(16), prvkeystd.Y.Text(16))
	}
	
	sx, sy, err := ecdsa.Sign(rand.Reader, prvkeytapp, rb)
	if err != nil{
		t.Fatal(err)
	}
	
	if !ecdsa.Verify(&prvkeystd.PublicKey, rb, sx, sy) {
		t.Fatal("verify signature fail (signatured by applied priv)")
	}
	
	sx, sy, err = ecdsa.Sign(rand.Reader, prvkeystd, rb)
	if err != nil{
		t.Fatal(err)
	}
	
	if !ecdsa.Verify(&prvkeytapp.PublicKey, rb, sx, sy) {
		t.Fatal("verify signature key fail (signatured by original priv)")
	}	
}

func TestDump_Privkey(t *testing.T){
	
	prvkeystd, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}
	
	prvkeyt := ECDSAPriv{1, prvkeystd.D}
	
	dmpstr, err := prvkeyt.DumpPrivKey()
	
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("dump privkey as", dmpstr)
	
	prvkeyresume, err := PrivKeyfromString(dmpstr)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if err != nil{
		t.Fatal(err)
	}
	
	if prvkeyresume.D.Cmp(prvkeystd.D) != 0 || prvkeyresume.CurveType != ECP256_FIPS186{
		t.Fatal("Unmatch private key number:", prvkeyresume.D.Text(16), prvkeystd.D.Text(16))		
	}
		
}
