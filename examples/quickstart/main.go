package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MetaRPC/GoCopier/copier"
)

func main() {
	fmt.Println("=== MetaRPC GoCopier Trade Replication Quick Start ===")
	ctx := context.Background()
	apiKey := "TRIAL"

	demo, err := copier.NewDemoAccountClient("https://mt5.mrpc.pro")
	if err != nil {
		log.Fatalf("Failed to initialize DemoAccountClient: %v", err)
	}

	var masterGuid, slaveGuid, copierId string

	svc, err := copier.NewCopierService("copy.mrpc.pro:443", apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize CopierService: %v", err)
	}
	defer svc.Close()

	defer func() {
		// 9. Cleanly Disconnect Terminal Sessions
		fmt.Println("\n[9] Disconnecting terminal sessions cleanly via /Disconnect...")
		if masterGuid != "" {
			discM, err := demo.Disconnect(ctx, masterGuid, apiKey)
			if err != nil {
				fmt.Printf("    Master disconnect error: %v\n", err)
			} else {
				fmt.Printf("    Master Terminal Cleanly Disconnected: %s (Lifetime: %ds)\n", discM.UniqueIdentifier, discM.FullLifeTimeSeconds)
			}
		}
		if slaveGuid != "" {
			discS, err := demo.Disconnect(ctx, slaveGuid, apiKey)
			if err != nil {
				fmt.Printf("    Slave disconnect error: %v\n", err)
			} else {
				fmt.Printf("    Slave Terminal Cleanly Disconnected:  %s (Lifetime: %ds)\n", discS.UniqueIdentifier, discS.FullLifeTimeSeconds)
			}
		}
		fmt.Println("\n=== GoCopier Trade Replication Completed Successfully ===")
	}()

	// 1. Provision live demo accounts
	fmt.Println("\n[1] Provisioning live demo accounts on MetaQuotes-Demo...")
	master, err := demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"}, apiKey)
	if err != nil {
		log.Fatalf("Failed to open master demo account: %v", err)
	}
	fmt.Printf("    Master Account Provisioned: #%d on %s\n", master.Login, master.Server)
	time.Sleep(1 * time.Second)

	slave, err := demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"}, apiKey)
	if err != nil {
		log.Fatalf("Failed to open slave demo account: %v", err)
	}
	fmt.Printf("    Slave Account Provisioned:  #%d on %s\n", slave.Login, slave.Server)
	time.Sleep(1 * time.Second)

	// 2. Connect terminals via ConnectEx with APIKey: TRIAL
	fmt.Printf("\n[2] Connecting terminals via ConnectEx (APIKey: %s)...\n", apiKey)
	for attempt := 1; attempt <= 3; attempt++ {
		connM, err := demo.ConnectEx(ctx, master.Login, master.Password, master.Server, apiKey)
		if err == nil && connM.TerminalInstanceGuid != "" {
			masterGuid = connM.TerminalInstanceGuid
			break
		}
		fmt.Printf("    Master ConnectEx attempt %d failed: %v. Retrying with fresh demo account...\n", attempt, err)
		master, _ = demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"}, apiKey)
		time.Sleep(2 * time.Second)
	}
	if masterGuid == "" {
		fmt.Println("    Failed to connect master terminal.")
		return
	}
	fmt.Printf("    Master Terminal Connected! GUID: %s\n", masterGuid)

	for attempt := 1; attempt <= 3; attempt++ {
		connS, err := demo.ConnectEx(ctx, slave.Login, slave.Password, slave.Server, apiKey)
		if err == nil && connS.TerminalInstanceGuid != "" {
			slaveGuid = connS.TerminalInstanceGuid
			break
		}
		fmt.Printf("    Slave ConnectEx attempt %d failed: %v. Retrying with fresh demo account...\n", attempt, err)
		slave, _ = demo.OpenDemoAccount(ctx, &copier.GuiDemoOpenAccountRequest{Server: "MetaQuotes-Demo"}, apiKey)
		time.Sleep(2 * time.Second)
	}
	if slaveGuid == "" {
		fmt.Println("    Failed to connect slave terminal.")
		return
	}
	fmt.Printf("    Slave Terminal Connected!  GUID: %s\n", slaveGuid)

	masterSessionId := copier.ToHyphenGuid(masterGuid)
	slaveSessionId := copier.ToHyphenGuid(slaveGuid)

	// 3. Start Trade Copier via gRPC on copy.mrpc.pro:443
	fmt.Println("\n[3] Starting Trade Copier via gRPC on copy.mrpc.pro:443...")
	startReq := &copier.StartRequest{
		UserKey:   apiKey,
		RiskType:  "LotMultiplier",
		RiskValue: "1.0",
		Master: &copier.Account{
			Type:     "MT5",
			User:     master.Login,
			Password: master.Password,
			Server:   master.Server,
			Id:       masterSessionId,
		},
		Slave: &copier.Account{
			Type:     "MT5",
			User:     slave.Login,
			Password: slave.Password,
			Server:   slave.Server,
			Id:       slaveSessionId,
		},
	}

	startRep, err := svc.Start(ctx, startReq)
	if err != nil {
		log.Fatalf("gRPC Start failed: %v", err)
	}
	fmt.Printf("    gRPC Start Reply: ok=%v, copierId=%s, error=%s\n", startRep.Ok, startRep.CopierId, startRep.Error)
	if !startRep.Ok {
		log.Fatalf("Copier Start returned error: %s", startRep.Error)
	}
	copierId = startRep.CopierId

	time.Sleep(4 * time.Second)

	// 4. Place Market Order on Master
	fmt.Println("\n[4] Opening Market Order on Master (0.01 EURUSD BUY)...")
	masterTicket, err := demo.OrderSend(ctx, masterGuid, "EURUSD", "TMT5_ORDER_TYPE_BUY", 0.01, apiKey)
	if err != nil {
		log.Fatalf("OrderSend failed: %v", err)
	}
	fmt.Printf("    Master Order Placed! Ticket: %d\n", masterTicket)

	// 5. Verify Trade Copied to Slave
	fmt.Println("\n[5] Verifying replicated trade on Slave account...")
	replicated := false
	for attempt := 1; attempt <= 15; attempt++ {
		time.Sleep(2 * time.Second)
		positions, err := demo.OpenedOrders(ctx, slaveGuid, apiKey)
		if err != nil {
			fmt.Printf("    Attempt %d error: %v\n", attempt, err)
			continue
		}
		fmt.Printf("    Attempt %d: Slave active positions count = %d\n", attempt, len(positions))
		if len(positions) > 0 {
			first := positions[0]
			fmt.Printf("    --> CONFIRMED ON SLAVE: Ticket=%d, Symbol=%s, Volume=%.2f, Type=%s\n", first.Ticket, first.Symbol, first.Volume, first.Type)
			replicated = true
			break
		}
	}

	if !replicated {
		fmt.Println("    WARNING: Slave trade replication timed out.")
	} else {
		fmt.Println("    SUCCESS: Trade successfully replicated to slave account!")
	}

	// 6. Close Position on Master
	if masterTicket != 0 {
		fmt.Printf("\n[6] Closing Master trade ticket #%d...\n", masterTicket)
		closeResp, err := demo.OrderClose(ctx, masterGuid, masterTicket, apiKey)
		if err != nil {
			fmt.Printf("    OrderClose error: %v\n", err)
		} else {
			fmt.Printf("    Master OrderClose result: %s\n", closeResp)
		}

		// 7. Verify Trade Closed on Slave
		fmt.Println("\n[7] Verifying trade closed on Slave...")
		for attempt := 1; attempt <= 15; attempt++ {
			time.Sleep(2 * time.Second)
			positions, _ := demo.OpenedOrders(ctx, slaveGuid, apiKey)
			if len(positions) == 0 {
				fmt.Println("    SUCCESS: Slave position closed by trade copier!")
				break
			}
			fmt.Printf("    Attempt %d: Slave positions still open: %d\n", attempt, len(positions))
		}
	}

	// 8. Remove Copier via gRPC
	if copierId != "" {
		fmt.Printf("\n[8] Removing Copier %s via gRPC...\n", copierId)
		remRep, err := svc.Remove(ctx, copierId)
		if err != nil {
			fmt.Printf("    Remove error: %v\n", err)
		} else {
			fmt.Printf("    Copier Remove Reply: ok=%v\n", remRep.Ok)
		}
	}
}

