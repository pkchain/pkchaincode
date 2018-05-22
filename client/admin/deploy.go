package main

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	
	"github.com/spf13/cobra"
	
	pb "gamecenter.mobi/paicode/protos"
	txutil "gamecenter.mobi/paicode/transactions"
)

var(
	isDebugMode *bool
	networkCode *int
	totalAmount *int64
	assignedDetail *[]string
)

func initDeployCmd() *cobra.Command{
	
	flags := deployCmd.PersistentFlags()
	
	isDebugMode = flags.Bool("debug", false, "use debug mode")
	networkCode = flags.Int("netcode", int(txutil.AddrHelper), "network code")
	totalAmount = flags.Int64("total", 1000000, "total amount of pais")
	assignedDetail = flags.StringSlice("assign", nil, "preassign pais to accounts")	
	
	deployCmd.AddCommand(deployOutCmd)
	
	return deployCmd
}

func cleanDeployFlags(){
	*assignedDetail = []string{}
}

func parseAssignDetail() ([]*pb.PreassignData, error){
	
	addrh := txutil.AddressHelper(*networkCode)
	data := make([]*pb.PreassignData, len(*assignedDetail))
	
	fail := func(i int, s string, err error) error{
		return errors.New(fmt.Sprintf("error at %dth string <%s>: %s", i, s, err))
	}
	
	for i, s := range *assignedDetail{
		ret := strings.Split(s, ":")
		
		if len(ret) != 2{
			return nil, fail(i, s, errors.New("Invalid preassign specification"))
		}
		
		amount, err := strconv.Atoi(ret[1])
		if err != nil{
			return nil, fail(i, s, err)
		}
		
		if b, err := addrh.VerifyUserId(ret[0]); !b{
			return nil, fail(i, s, errors.New(fmt.Sprint("Invalid address:", err)))
		}
		
		data[i] = &pb.PreassignData{Userid: ret[0], Pais: int64(amount)}
	}
	
	return data, nil
}


var deployCmd = &cobra.Command{
	Use:   "deploy [options...]",
	Short: fmt.Sprintf("deploy chaincode to fabric"),	
	
}

var deployOutCmd = &cobra.Command{
	Use:   "outprint",
	Short: fmt.Sprintf("simply print out the ctormsg"),
	
	RunE: func(cmd *cobra.Command, args []string) error {
		
		preassign, err := parseAssignDetail()
		if err != nil{
			return err
		}
		
		initset := &pb.InitChaincode{&pb.DeploySetting{*isDebugMode, int32(*networkCode), *totalAmount, *totalAmount}, preassign}
		arg, err := txutil.EncodeChaincodeTx(initset)
		if err != nil {
			return err
		}
		
		fmt.Println("Init msg:", arg)
		return nil		
	},	
}



