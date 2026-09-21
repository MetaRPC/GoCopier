package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MetaRPC/GoCopier/copier"
)

func main() {
	fmt.Println("=== GoCopier Quick Start Demo ===")
	ctx := context.Background()
	apiKey := "TRIAL"

	demo, err := copier.NewDemoAccountClient("https://mt5.mrpc.pro")
	if err != nil {
		log.Fatalf("Failed to initialize DemoAccountClient: %v", err)
	}

	// 1. Provision live demo account
	fmt.Println("\n[1] Provisioning live demo account on MetaQuotes-Demo...")
	master, err := demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"}, apiKey)
	if err != nil {
		log.Fatalf("Failed to open demo account: %v", err)
	}
	fmt.Printf("    Master Account Provisioned: #%d on %s\n", master.Login, master.Server)

	// 2. Connect terminal via ConnectEx with APIKey: TRIAL
	fmt.Printf("\n[2] Connecting terminal via ConnectEx (APIKey: %s)...\n", apiKey)
	conn, err := demo.ConnectEx(ctx, master.Login, master.Password, master.Server, apiKey)
	if err != nil {
		log.Fatalf("ConnectEx failed: %v", err)
	}
	fmt.Printf("    Terminal Connected! Instance GUID: %s\n", conn.TerminalInstanceGuid)

	// 3. Interacting with Copier Service
	fmt.Printf("\n[3] Interacting with Copier Service (user_key: %s)...\n", apiKey)
	svc, err := copier.NewCopierService("copy.mrpc.pro:443", apiKey)
	if err == nil {
		defer svc.Close()
		listReply, _ := svc.List(ctx)
		if listReply != nil {
			fmt.Printf("    Active copiers for %s: %d\n", apiKey, len(listReply.Copiers))
		}
	}

	// 4. Cleanly Disconnect Terminal Session
	fmt.Printf("\n[4] Disconnecting terminal session %s...\n", conn.TerminalInstanceGuid)
	disc, err := demo.Disconnect(ctx, conn.TerminalInstanceGuid, apiKey)
	if err != nil {
		log.Fatalf("Disconnect failed: %v", err)
	}
	fmt.Printf("    Terminal Cleanly Disconnected: %s (Lifetime: %ds)\n", disc.UniqueIdentifier, disc.FullLifeTimeSeconds)

	fmt.Println("\n=== GoCopier Quick Start Completed Successfully ===")
}
