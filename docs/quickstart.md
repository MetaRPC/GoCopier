# First request with Go

List your own copiers to check connectivity and identity without creating a terminal or placing a trade.

## 1. Install

Requires **Go 1.21+**. These instructions build from source; they do not assume a registry release exists.

```bash
git clone https://github.com/MetaRPC/GoCopier.git
cd GoCopier
go mod download
go test ./...
```

Java on Windows uses `;` for classpath separators. Run C++ builds in a Visual Studio developer shell. Activate your Python virtual environment before installing.

## 2. Set your owner key

=== "Bash"
    ```bash
    export MRPC_API_KEY='your-owner-key'
    ```
=== "PowerShell"
    ```powershell
    $env:MRPC_API_KEY = 'your-owner-key'
    ```

Use the real owner API key, not a broker password, terminal ID or `TRIAL`. Read [authentication](authentication.md) if access is refused.

## 3. Run

```bash
go run ./examples/list
```

Source: [`examples/list/main.go`](https://github.com/MetaRPC/GoCopier/blob/main/examples/list/main.go).

```go
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
```

A successful empty list means this owner has no copiers. Treat failed requests as errors, not empty accounts.

## 4. Prepare a demo workflow

Read [the lifecycle guide](lifecycle.md), then [choose sizing](sizing.md). Existing replication examples under `examples/` are advanced integration examples: they may provision accounts and trade. Inspect them and their cleanup before running with disposable demos. List is the default starting point.
