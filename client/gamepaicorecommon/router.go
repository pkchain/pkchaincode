package gamepaicorecommon

import (
	"fmt"
	"github.com/gocraft/web"
	clicore "gamecenter.mobi/paicode/client"
)

var DefClient *clicore.ClientCore

type GamepaiREST struct{
	
}

type AccountREST struct{
	*GamepaiREST
	id string
	shouldPersist bool
}

type RpcQueryREST struct{
	*GamepaiREST
	workCore *clicore.RpcCore 
}

type RpcREST struct{
	*GamepaiREST
	RpcQueryREST
}

type BlockUtilREST struct{
	*GamepaiREST
}

type restData struct{
	Status string `json:"status"`
	Data   interface{} `json:"data"`
}


func (s *GamepaiREST) SetResponseType(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS (default option handler will handle OPTION and set Access-Control-Allow-Method properly)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type")

	next(rw, req)
}

func (s *AccountREST) PersistAccount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	next(rw, req)
	
	if s.shouldPersist {
		err := DefClient.Accounts.KeyMgr.Persist()
		
		if err != nil{
			panic(fmt.Sprintln("Persist fail", err))
		}
		
	} 
}

func (s *RpcQueryREST) LoadRPC(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	
	s.workCore = clicore.RpcCoreFromClient(&DefClient.Rpc)
	
	err := s.workCore.Rpc.Rpcbuilder.VerifyConn()
	if err != nil{
		s.RpcFail(rw, req, err.Error())
		return
	}		
	
	next(rw, req)
}

func (s *RpcREST) SetRPCKey(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	err := req.ParseForm()
	if err != nil{
		panic(err)
	}

	kid := req.Form.Get("id")
	if kid == "" {
		panic("Must specific id")
	}
	
	key, err := DefClient.Accounts.KeyMgr.LoadPrivKey(kid)
	if err != nil{
		panic(fmt.Sprintf("No corresponding privkey for %s", kid))
	}		
	
	s.workCore.Rpc.PrivKey = key
	next(rw, req)
}

func BuildRouter() *web.Router {
	
	router := web.New(GamepaiREST{})

	// Add middleware
	router.Middleware((*GamepaiREST).SetResponseType) 

	accRouter := router.Subrouter(AccountREST{shouldPersist:false}, "/account")
	accRouter.Middleware((*AccountREST).PersistAccount)
	accRouter.Post("/", (*AccountREST).NewAcc)
	accRouter.Get("/", (*AccountREST).ListAcc)
	accRouter.Delete("/:id", (*AccountREST).DeleteAcc)
	accRouter.Get("/:id", (*AccountREST).QueryAcc)	
	accRouter.Get("/dump/:id", (*AccountREST).DumpAcc)
	accRouter.Get("/verify/:addr", (*AccountREST).VerifyAddr)
	
	rpcRouter := router.Subrouter(RpcREST{}, "/rpc")
	rpcRouter.Middleware((*RpcREST).LoadRPC)
	rpcRouter.Middleware((*RpcREST).SetRPCKey)
	rpcRouter.Post("/registar", (*RpcREST).Registar)
	rpcRouter.Post("/fund", (*RpcREST).Fund)

	chainRouter := router.Subrouter(RpcQueryREST{}, "/chain")
	chainRouter.Middleware((*RpcQueryREST).LoadRPC)
	chainRouter.Get("/:addr", (*RpcQueryREST).QueryUser) 
	chainRouter.Get("/", (*RpcQueryREST).QueryChain)

	blockRouter := router.Subrouter(BlockUtilREST{}, "/blocks")
	blockRouter.Post("/parsetxpayload", (*BlockUtilREST).ParsePayload)

	return router
}

