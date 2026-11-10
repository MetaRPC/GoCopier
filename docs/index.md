---
hide:
  - toc
---
<div class="hero"><span class="eyebrow">METARPC / DEVELOPER GUIDES</span><h1>Trade copying.<br>Built with Go.</h1><p>From your first authenticated request to a monitored master–follower workflow. Clear examples, explicit sizing and practical recovery.</p><a class="md-button md-button--primary" href="quickstart/">Start building →</a> <a class="md-button" href="reference/">Explore the API</a><div class="hero-meta">TLS transport · Owner authentication · Demo-first validation</div></div>

<div class="grid cards" markdown>

- **01 · Connect with confidence**

    Understand keys, account IDs and transport before sending a request.

    [Authentication →](authentication.md)

- **02 · Control the lifecycle**

    Start once, inspect state, pause and clean up your test.

    [Copier lifecycle →](lifecycle.md)

- **03 · Make sizing explicit**

    Choose lot, balance or equity sizing and account for broker constraints.

    [Sizing and settings →](sizing.md)

- **04 · Diagnose the result**

    Separate network success from management success and broker execution.

    [Troubleshooting →](troubleshooting.md)

</div>

## Your first successful request

Requires Go 1.21+. Install dependencies and run the read-only List example. Keep your key in `MRPC_API_KEY`.

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

[Full setup](quickstart.md) · [Testing guide](testing.md) · [All SDKs](sdks.md)
