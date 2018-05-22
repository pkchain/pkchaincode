package gamepaicorecommon

import (
	"github.com/gocraft/web"
	"net/http"
	"encoding/json"
	_ "gamecenter.mobi/paicode/client"
)

func (s *RpcREST) Registar(rw web.ResponseWriter, req *web.Request){
	txid, err := s.workCore.Rpc.Registry()

	encoder := json.NewEncoder(rw)
	if err != nil{
		rw.WriteHeader(http.StatusServiceUnavailable)
		encoder.Encode(restData{"tx fail", err.Error()})
		return
	}	
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", txid})	
}

func (s *RpcREST) Fund(rw web.ResponseWriter, req *web.Request){

	if len(req.Form) == 0{
		panic("form must be parsed before")
	}
	
	toaddr := req.Form.Get("to")
	amount := req.Form.Get("amount")
	msg := req.Form.Get("message")
	
	var txid string
	var err error
	
	if msg == ""{
		txid, err = s.workCore.Rpc.Fund(toaddr, amount)
	}else{
		txid, err = s.workCore.Rpc.Fund(toaddr, amount, msg)
	}
	
	encoder := json.NewEncoder(rw)
	if err != nil{
		rw.WriteHeader(http.StatusServiceUnavailable)
		encoder.Encode(restData{"tx fail", err.Error()})
		return
	}	
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", txid})
	
}