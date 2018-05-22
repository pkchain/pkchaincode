package security

import(
	
	"strings"
	"github.com/op/go-logging"
	"github.com/hyperledger/fabric/core/chaincode/shim"	
	"github.com/hyperledger/fabric/core/chaincode/shim/crypto/attr"	
	pb "gamecenter.mobi/paicode/protos" 
)

type SecurityPolicy struct{
	dbgMode bool
	netCode int
}

const(
		
	AdminPrivilege string = "Admin"
	ManagerPrivilege string = "Manager"
	DelegatePrivilege string = "Delegate"
	ObserverPrivilege string = "Observer"
	
	debugPrivilege string = "debug"
	noPrivilege string = "none"
	
	Privilege_Attr string = "PaiAdminRole"
	Region_Attr string = "PaiAdminRegion"
	
	debugRegion string = "debug"
	noRegion string = "none"
)

var logger = logging.MustGetLogger("chaincode_sec")

//keep a singleton
var Helper = &SecurityPolicy{true, 0}

func (s *SecurityPolicy) Update(set *pb.DeploySetting){
	s.dbgMode = set.DebugMode
	s.netCode = int(set.NetworkCode)
}

func InitSecHelper(set *pb.DeploySetting){	
	Helper.Update(set)
}

func (sec *SecurityPolicy) ActiveAudit(stub shim.ChaincodeStubInterface, desc string){
	
}

func (sec *SecurityPolicy) GetPrivilege(stub shim.ChaincodeStubInterface) (string, string){
		
	attrHandler, err := attr.NewAttributesHandlerImpl(stub)
	if err != nil{
		logger.Error("Create Attr handler fail", err)
		
		if sec.dbgMode{
			return debugPrivilege, debugRegion
		}		
		
		return noPrivilege, noRegion
	}

	var privstr, regionstr string
	privilege, err := attrHandler.GetValue(Privilege_Attr)
	if err != nil{
		logger.Error("get privilege attr fail", err)
		privstr = noPrivilege
	}else{
		privstr = string(privilege)
	}
	
	region, err := attrHandler.GetValue(Region_Attr)
	if err != nil{
		logger.Error("get region attr fail", err)
		regionstr = noRegion
	}else{
		regionstr = string(region)
	}
	
	if sec.dbgMode{
		return debugPrivilege, debugRegion
	}		
	return privstr, regionstr
	
}

func (sec *SecurityPolicy) VerifyPrivilege(certpriv string, expect string) bool{
	
	if sec.dbgMode && strings.Compare(certpriv, debugPrivilege) == 0{
		return true
	}
	
	if strings.Compare(certpriv, expect) != 0{
		logger.Error("Comprivilege fail: ", certpriv, expect)
		return false
	}
	
	return true
}

func (sec *SecurityPolicy) VerifyRegion(region string, expect string) bool{
	
	if sec.dbgMode && strings.Compare(region, debugRegion) == 0{
		return true
	}	
	
	return strings.Compare(region, expect) == 0
	
}
