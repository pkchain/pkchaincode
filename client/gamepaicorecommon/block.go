package gamepaicorecommon

import (
	"github.com/gocraft/web"
	"net/http"
	"encoding/json"	
)

func (_ *BlockUtilREST) ParsePayload(rw web.ResponseWriter, req *web.Request){
	
	err := req.ParseForm()
	if err != nil{
		panic(err)
	}	
		
	encoder := json.NewEncoder(rw)	
	out, err := DefClient.Block.DecodePayload(req.Form.Get("payload"))
	
	if err != nil{
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restData{"Invalid payload", err.Error()})
	}else{
		encoder.Encode(restData{"ok", out})
	}
}