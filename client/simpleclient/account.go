package main

import (
	"fmt"
	"errors"
	
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account [command...]",
	Short: fmt.Sprintf("account commands."),
}

var listPrivkeyCmd = &cobra.Command{
	Use:       "list",
	Short:     fmt.Sprintf("list all keys"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		ret := defClient.Accounts.ListKeyData(args...)
		
		fmt.Println("-----------------------------------------------------------------")
		fmt.Println("|  No. |       Name         |               Address             |")
		fmt.Println("-----------------------------------------------------------------")
		
		for i, v := range ret{
			fmt.Printf("|%6d|%20s|%35s|\n", i, v[0], v[1])
		}
		
		fmt.Println("-----------------------------------------------------------------")
		
		return nil
	},
}

var queryPrivkeyCmd = &cobra.Command{
	Use:       "query <remark>",
	Short:     fmt.Sprintf("query a key"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		addr, err := defClient.Accounts.GetAddress(args...)
		if err != nil{
			return err
		}
		
		fmt.Println(args[0], ":", addr)
		
		return nil
	},
}

var genPrivkeyCmd = &cobra.Command{
	Use:       "generate [remark]",
	Short:     fmt.Sprintf("generate a privkey"),
	Long:      fmt.Sprintf(`generate a privkey and save it with the name of remark.`),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		rmk, err := defClient.Accounts.GenPrivkey(args...)
		if err != nil{
			return err
		}
		
		if len(args) < 1{
			fmt.Println("Done:", rmk)
		}else{
			fmt.Println("Done")
		}
		return nil
	},
}

var dumpPrivkeyCmd = &cobra.Command{
	Use:       "dump <remark>",
	Short:     fmt.Sprintf("dump out a privkey"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		ret, err := defClient.Accounts.DumpPrivkey(args...)
		
		if err != nil{
			return err
		}
		
		fmt.Println(ret)		
		return nil	
	},
}

var importPrivkeyCmd = &cobra.Command{
	Use:       "import <dump string> [remark]",
	Short:     fmt.Sprintf("Import a privkey"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		ret, err := defClient.Accounts.ImportPrivkey(args...)
		
		if err != nil{
			return err
		}
		
		if len(args) < 1{
			fmt.Println("Import done as", ret)
		}else{
			fmt.Println("Import done")
		}	
		return nil	
	},
}

var usePrivkeyCmd = &cobra.Command{
	Use:       "use <remark>",
	Short:     fmt.Sprintf("Use a privkey for rpc calling"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		if len(args) != 1{
			return errors.New("Invalid argument count")
		}
		
		key, err := defClient.Accounts.KeyMgr.LoadPrivKey(args[0])
		if err != nil{
			return err
		}
		
		defClient.Rpc.PrivKey = key	
		fmt.Println("Done")
		return nil	
	},
}


func init(){
	accountCmd.AddCommand(genPrivkeyCmd)
	accountCmd.AddCommand(dumpPrivkeyCmd)
	accountCmd.AddCommand(listPrivkeyCmd)
	accountCmd.AddCommand(queryPrivkeyCmd)
	accountCmd.AddCommand(importPrivkeyCmd)
	accountCmd.AddCommand(usePrivkeyCmd)
}