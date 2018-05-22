package client

import (
    "testing"
    "encoding/json"
)

func TestBlockExploere(t *testing.T) {
	be := blockExplorer{}
	
	ret, err := be.DecodePayload("CpQDCAESERIPZ2FtZXBhaWNvcmVfdjAxGt4CCglVU0VSX0ZVTkQKtAFDaUpCV2pKT04zcFVSRFpvVEhwaFJWTXpVbFJVYW00MU9FOUhOMk5JZDBkTGJYUm5FbUhtczZqbWhJL3Z2SnJvcjdma3VJM29wb0hvdjU3bnU2MHk1cXloNW9pVzVhU2E1cXloNVpDUjU1dTQ1WkNNNVp5dzVaMkE1WStSNllDQjVaQ001cUMzNXBXdzZZZVA1WktNNVpDTTVxQzM1YVNINXJPbzU1cUU2TDJzNkxTbTQ0Q0MKOENpWUlDUklpUVZWWFVqQmFYekZsTjBsek5EUjRhVk40TjBRNWJuazVkVWwwY0RWT1UyaGpRUT09CmBDa1FLSURXM012RUVQbllYdWhIazZlVzFEVHRVYVh1MEtvWTBSNzQ0Mkhwa0laSVpFaUE2bzdwYzZLS0I5MTBoREZFZ2pRUVo0NUxRa2YwTDVBcWVFZWsyY2czQWVRPT1CDFBhaUFkbWluUm9sZUIOUGFpQWRtaW5SZWdpb24=")
	
	if err != nil{
		t.Fatal(err)
	}
	
	js, err := json.MarshalIndent(ret, "", " ")
	if err != nil{
		t.Fatal(err)
	}
		
	t.Log(string(js))
}

