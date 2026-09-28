# GoCopier

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Docs](https://img.shields.io/badge/docs-online-green.svg)](https://github.com/MetaRPC/GoCopier/tree/main/docs)

Official Go SDK for the MetaRPC Trade Copier high-performance trade replication engine via gRPC (`copy.mrpc.pro:443`).

## Installation

```bash
go get github.com/MetaRPC/GoCopier
```

## 🏃 How to Run Examples

Clone the repository and run the trade copier example out-of-the-box:

```bash
git clone https://github.com/MetaRPC/GoCopier.git
cd GoCopier/examples/quickstart
go mod download

# 1. Run with default TRIAL key:
go run main.go

# 2. Or pass your MetaRPC API key directly as an argument:
go run main.go your_api_key_here

# 3. Or use the MRPC_API_KEY environment variable:
export MRPC_API_KEY="your_api_key_here"        # Windows CMD: set MRPC_API_KEY=your_api_key_here
go run main.go                                 # Windows PowerShell: $env:MRPC_API_KEY="your_api_key_here"
```

## Quick Start

See [Quick Start Documentation](https://github.com/MetaRPC/GoCopier/blob/main/docs/All_Guides/Your_First_Project.md) for a 10-minute walkthrough.
