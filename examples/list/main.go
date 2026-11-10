package main
import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/MetaRPC/GoCopier/copier"
)
func main() {
    service, err := copier.NewCopierService("copy.mrpc.pro:443", os.Getenv("MRPC_API_KEY"))
    if err != nil { log.Fatal(err) }
    defer service.Close()
    result, err := service.List(context.Background())
    if err != nil { log.Fatal(err) }
    if !result.Ok { log.Fatal(result.Error) }
    for _, item := range result.Copiers { fmt.Println(item.Id, item.Paused) }
}
