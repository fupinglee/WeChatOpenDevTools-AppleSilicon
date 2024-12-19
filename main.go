package main

import (
	"fmt"
	"os"

	"WeChatOpenDevTools-AppleSilicon/pkg/utils"

	"github.com/spf13/cobra"
)

func main() {
	var pid int
	commons := utils.NewCommons()

	rootCmd := &cobra.Command{
		Use:   "WeChatOpenDevTools-AppleSilicon",
		Short: "WeChat DevTools for Apple Silicon",
		Long:  `WeChat DevTools - A tool for injecting and hooking WeChat process`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if pid != 0 {
				err = commons.LoadWeChatConfigsPID(pid)
			} else {
				err = commons.LoadWeChatConfigs()
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().IntVarP(&pid, "pid", "p", 0, "Specify WeChat process PID")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
