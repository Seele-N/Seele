package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/seele-n/seele/cmd/seeled/cmd"
    svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/seele-n/seele/app"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
    if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
