package client

import (
	"gamecenter.mobi/paicode/wallet"
	sec "gamecenter.mobi/paicode/chaincode/security"
	"github.com/hyperledger/fabric/peerex"
)

type ClientCore struct{
	Accounts accountManager
	Rpc		 rpcManager
	Block    blockExplorer
}

type RpcCore struct{
	Rpc		 rpcManager
}

func NewClientCore(config *peerex.GlobalConfig) *ClientCore{
	
	walletmgr := wallet.CreateSimpleManager(config.GetPeerFS() + "wallet.dat")
	
	return &ClientCore{Accounts: accountManager{walletmgr}}
}

func RpcCoreFromClient(rpc *rpcManager) *RpcCore{
	
	c := new(RpcCore)
	
	c.Rpc.Rpcbuilder = &peerex.RpcBuilder{}
	c.Rpc.Rpcbuilder.Conn = rpc.Rpcbuilder.Conn
	c.Rpc.Rpcbuilder.ChaincodeName = rpc.Rpcbuilder.ChaincodeName
	c.Rpc.PrivKey = rpc.PrivKey
	c.Rpc.Rpcbuilder.Security = rpc.Rpcbuilder.Security
	return c
}

func (c *ClientCore) IsRpcReady() bool{
	return c.Rpc.Rpcbuilder != nil
}


func (c *ClientCore) PrepareRpc(conn peerex.ClientConn){
	c.Rpc.Rpcbuilder = &peerex.RpcBuilder{}
	c.Rpc.Rpcbuilder.Conn = conn
}

func (c *ClientCore) SetRpcRegion(user string){
	if c.Rpc.Rpcbuilder == nil{
		panic("Must call PrepareRpc first")
	}
	c.Rpc.Rpcbuilder.Security = &peerex.SecurityPolicy{User: user, 
		Attributes: []string{sec.Privilege_Attr, sec.Region_Attr}}
}

func (c *ClientCore) ReleaseRpc(){
	
	if c.Rpc.Rpcbuilder != nil && c.Rpc.Rpcbuilder.Conn.C != nil{
		c.Rpc.Rpcbuilder.Conn.C.Close()	
	}
	
}


