package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MetaRPC/GoCopier/copier"
)

func main() {
	ctx := context.Background()
	demo, _ := copier.NewDemoAccountClient("mt5.mrpc.pro:443")
	master, _ := demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"})
	slave, _ := demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"})
	fmt.Printf("Master: %d, Slave: %d\n", master.Login, slave.Login)

	svc, err := copier.NewCopierService("copy.mrpc.pro:443", "YOUR_USER_KEY", "YOUR_MANAGER_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()

	reply, _ := svc.Start(ctx, &copier.StartRequest{
		Master:    &copier.Account{Type: "MT5", User: master.Login, Password: master.Password, Server: master.Server},
		Slave:     &copier.Account{Type: "MT5", User: slave.Login, Password: slave.Password, Server: slave.Server},
		RiskType:  "LotMultiplier",
		RiskValue: "1.5",
	})
	fmt.Printf("Started Copier: %s\n", reply.CopierId)
}
