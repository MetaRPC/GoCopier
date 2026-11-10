// Native SDK adapter. Credentials are read from stdin.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MetaRPC/GoCopier/copier"
	"os"
	"strconv"
)

func main() {
	scan := bufio.NewScanner(os.Stdin)
	v := make([]string, 13)
	for i := range v {
		if !scan.Scan() {
			panic("Expected 13 configuration lines")
		}
		v[i] = scan.Text()
	}
	c, err := copier.NewCopierService(v[0], v[1])
	if err != nil {
		json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": false, "error": err.Error()})
		return
	}
	defer c.Close()
	account := func(i int) *copier.Account {
		u, e := strconv.ParseUint(v[i], 10, 64)
		if e != nil {
			panic(e)
		}
		return &copier.Account{Type: "MT5", User: u, Password: v[i+1], Server: v[i+2], Id: v[i+3]}
	}
	var r any
	ctx := context.Background()
	switch os.Args[1] {
	case "start":
		r, err = c.Start(ctx, &copier.StartRequest{Master: account(2), Slave: account(6), RiskType: v[10], RiskValue: v[11], FixedMasterBalance: v[12], CopySl: true, CopyTp: true})
	case "list":
		r, err = c.List(ctx)
	case "pause":
		r, err = c.Pause(ctx, os.Args[2], os.Args[3] == "true")
	case "remove":
		r, err = c.Remove(ctx, os.Args[2])
	default:
		err = fmt.Errorf("Unknown operation")
	}
	if err != nil {
		r = map[string]any{"ok": false, "error": err.Error()}
	}
	json.NewEncoder(os.Stdout).Encode(r)
}
