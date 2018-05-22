package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	
	"github.com/hyperledger/fabric/peerex"
)

var logger = peerex.InitLogger("main")

func main() {
	
	config := peerex.GlobalConfig{}
	err := config.InitGlobal(true)
	
	if err != nil{
		panic(err)		
	}
	
	var default_conn peerex.ClientConn
	err = default_conn.Dialdefault()
	if err != nil{
		fmt.Println("Dial to peer fail:", err)
		os.Exit(1)
	}	
	
	defer default_conn.C.Close()
	
	fmt.Print("Input chaincode name: ")
	var chaincodename string = ""
	fmt.Scan(&chaincodename)
	if len(chaincodename) == 0{
		fmt.Println("Invalidate chaincode name")
		os.Exit(1)
	}
	fmt.Println("Chaincode name is:", chaincodename)

	fmt.Print("Now input function name: ")
	var funcname string = ""
	fmt.Scan(&funcname)
	if len(funcname) == 0{
		fmt.Println("Invalidate function name")
		os.Exit(1)
	}
	fmt.Println("Function name is:", funcname)
	
	fmt.Println("You can start calling the function:")
	
	reader := bufio.NewReader(os.Stdin)	
	var ln string = "" 
	
	for {
		retbyte, notfinished, err := reader.ReadLine()
		if err != nil{
			break
		}
		
		ln += string(retbyte)
		if !notfinished {
			//handle read command line
			//fmt.Println("We get:", ln)
			args := strings.Split(ln, " ")
			if args == nil{
				fmt.Println("Must input some arguments")		
			}else{
				logger.Debug(len(args), "items:", args)
				
				rpc := peerex.RpcBuilder{ChaincodeName: chaincodename, Function: funcname, Conn: default_conn}
				_, err := rpc.Fire(args)
				if err != nil{
					fmt.Println("Command error:", err)		
				}				
			}
			
			ln = ""
			fmt.Println("Continue to next rpc:")
		}
	}
	
	fmt.Println("Exiting ...")
	
}

