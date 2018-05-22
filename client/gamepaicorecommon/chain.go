package gamepaicorecommon

import (
	
	"net/http"
	"github.com/gocraft/web"
	"encoding/json"
	_ "gamecenter.mobi/paicode/client"
)

func reconstructRpcRet(jsonstr []byte) (*restData, error){
	
	var outobj interface{}
	err := json.Unmarshal(jsonstr, &outobj)
	
	if err != nil{
		return nil, err
	}
	
	return &restData{"ok", outobj.(map[string]interface{})}, nil
	
}

func (s *RpcQueryREST) RpcFail(rw web.ResponseWriter, req *web.Request, reason string){
	encoder := json.NewEncoder(rw)	
	rw.WriteHeader(http.StatusServiceUnavailable)
	encoder.Encode(restData{"network broken", reason})		
}

func (s *RpcQueryREST) QueryUser(rw web.ResponseWriter, req *web.Request){
	address := req.PathParams["addr"]
	if address == "" {
		panic("Must specific address")
	}
	
	encoder := json.NewEncoder(rw)
	ret, err := s.workCore.Rpc.QueryUser(address)
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		encoder.Encode(restData{"Query fail", err.Error()})
		return
	}
	
	data, err := reconstructRpcRet(ret)
	if err != nil{
		panic(err)
	}
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(data)
}



func (s *RpcQueryREST) QueryChain(rw web.ResponseWriter, req *web.Request){
	encoder := json.NewEncoder(rw)
	ret, err := s.workCore.Rpc.QueryGlobal()
	if err != nil{
		rw.WriteHeader(http.StatusServiceUnavailable)
		encoder.Encode(restData{"Query fail", err.Error()})
		return
	}

	data, err := reconstructRpcRet(ret)
	if err != nil{
		panic(err)
	}
	
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(data)	
}

