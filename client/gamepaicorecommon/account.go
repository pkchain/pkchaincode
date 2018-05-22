package gamepaicorecommon

import (
	
	_ "fmt"
	"net/http"
	"github.com/gocraft/web"
	"encoding/json"
	tx "gamecenter.mobi/paicode/transactions"
	_ "gamecenter.mobi/paicode/client"
)

func (s *AccountREST) QueryAcc(rw web.ResponseWriter, req *web.Request){

	id := req.PathParams["id"]
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Account id not found\"}"))//just write a raw json
		return
	}
	
	encoder := json.NewEncoder(rw)
	addr, err := DefClient.Accounts.GetAddress(id)
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		encoder.Encode(restData{"Account not exist", err.Error()})
		return
	}
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", addr})
	
}

func (s *AccountREST) ListAcc(rw web.ResponseWriter, req *web.Request){
	
	retmap := map[string]string{}
	for _, v := range DefClient.Accounts.ListKeyData(){
		retmap[v[0]] = v[1]
	}
	
	encoder := json.NewEncoder(rw)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", retmap})
}

func (s *AccountREST) DeleteAcc(rw web.ResponseWriter, req *web.Request){
	id := req.PathParams["id"]
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Account id not found\"}"))//just write a raw json
		return
	}
	
	encoder := json.NewEncoder(rw)
	err := DefClient.Accounts.DelPrivkey(id)
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		encoder.Encode(restData{"Account not exist", nil})
		return
	}
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", nil})
	s.shouldPersist = true
}

func (s *AccountREST) NewAcc(rw web.ResponseWriter, req *web.Request){
	
	err := req.ParseForm()
	encoder := json.NewEncoder(rw)
	
	if err != nil || len(req.Form) == 0{
		rw.WriteHeader(http.StatusNotAcceptable)
		encoder.Encode(restData{"Wrong form or not expected content (application/x-www-urlencoded)", err.Error()})
		return		
	}
		
	accountid := req.Form["id"]	
	if len(accountid) == 0{
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Must provide account id\"}"))
		return				
	}
	
	prvkstr := req.Form["privatekey"]
		
	if len(prvkstr) != 0{
		//import
		_, err = DefClient.Accounts.ImportPrivkey(prvkstr[0], accountid[0])
	}else{
		//generate
		_, err = DefClient.Accounts.GenPrivkey(accountid[0])
	}
	
	if err != nil{
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restData{"Import fail", err.Error()})
		return
	}
	
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{\"status\":\"ok\"}"))
	s.shouldPersist = true
}

func (s *AccountREST) DumpAcc(rw web.ResponseWriter, req *web.Request){
	
	id := req.PathParams["id"]
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Account id not found\"}"))//just write a raw json
		return
	}	
	
	encoder := json.NewEncoder(rw)
	ret, err := DefClient.Accounts.DumpPrivkey(id)

	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		encoder.Encode(restData{"Account not exist", err.Error()})
		return
	}	
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", ret})
		
}

func (s *AccountREST) VerifyAddr(rw web.ResponseWriter, req *web.Request){

	addr := req.PathParams["addr"]
	if addr == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Address not found\"}"))//just write a raw json
		return
	}

	encoder := json.NewEncoder(rw)
	ret, _ := tx.AddrHelper.VerifyUserId(addr)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restData{"ok", ret})	
	
}
