package main 

import (
	"os"
	"fmt"
	"bufio"
	_ "strings"
	
	arghelper "github.com/mattn/go-shellwords"
	
	"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/peerex"
)

var mainCmd = &cobra.Command{
	Use: ">",
}

func main() {
	
	config := peerex.GlobalConfig{}
	err := config.InitGlobal(true)
	
	if err != nil{
		panic(err)		
	}
	
	fmt.Print("Starting .... ")
	mainCmd.AddCommand(initDeployCmd())
	
	mainCmd.SetArgs([]string{"help"})
	err = mainCmd.Execute()
	if err != nil{
		fmt.Println("Command handler error:", err)
		os.Exit(1)		
	}	
	
	
	reader := bufio.NewReader(os.Stdin)
	parser := arghelper.NewParser()
	
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
			args, err := parser.Parse(ln)
			if err != nil{
				fmt.Println("Input parse error:", err)		
			}else{
				mainCmd.SetArgs(args)
				err = mainCmd.Execute()
				if err != nil{
					fmt.Println("Command error:", err)		
				}				
			}
			
			ln = ""
			cleanDeployFlags()
			fmt.Println("Continue to next command:")
		}
	}
	
	fmt.Println("Exiting ...")	
	
}

