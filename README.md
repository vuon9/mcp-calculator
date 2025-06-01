# MCP Calculator

A demonstration of a Model Context Protocol (MCP) server built using the `github.com/mark3labs/mcp-go` library. It exposes a simple "calculate" tool.

## Features

*   **MCP Server Implementation**: Built on `github.com/mark3labs/mcp-go/server`, extended from the default example server.
*   **Calculator Tool**:
    *   Name: `calculate`
    *   Description: Perform basic arithmetic operations.
    *   Operations: `add`, `subtract`, `multiply`, `divide`.
    *   Accepts an array of numbers for calculations.
*   **HTTP Transport**: Uses `StreamableHTTPServer` for communication.

## Prerequisites

*   Go

## How to Run

1.  **Clone:**
    ```bash
    git clone https://github.com/vuon9/mcp-calculator.git
    cd mcp-calculator
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the server:**
    ```bash
    go run main.go
    ```

4.  The server will start and listen on `http://localhost:8080` and the default MCP endpoint will be `http://localhost:8080/mcp`.

## How to Use

The server communicates using the Model Context Protocol over HTTP. You'll typically interact with it by sending JSON-RPC messages to the `/mcp` endpoint.

### 1. Initialize Session

First, you need to initialize a session. The server is stateful by default and will return a session ID.

**Request (POST to `http://localhost:8080/mcp`):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "clientInfo": {
      "name": "my-calculator-client",
      "version": "0.1.0"
    }
  }
}
```

**Example Response:**
The server will respond with its capabilities and, importantly, an `Mcp-Session-Id` header if it's stateful (which it is by default).
```
HTTP/1.1 200 OK
Content-Type: application/json
Mcp-Session-Id: <session-id-from-server>

{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2025-03-26",
    "serverInfo": {
      "name": "Calculator Demo",
      "version": "1.0.0"
    },
    "capabilities": {
      // ... capabilities including tools ...
    }
  }
}
```
**Note:** Store the `Mcp-Session-Id` from the response header for subsequent requests.

### 2. Call the Calculator Tool

To use the calculator:

**Request (POST to `http://localhost:8080/mcp`):**
*Include the `Mcp-Session-Id` header.*

**Example: Add 10, 20, and 5**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "calculate",
    "arguments": {
      "operation": "add", // or "subtract", "multiply", "divide"
      "numbers": [10, 20, 5]
    }
  }
}
```
**Expected Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "35.00"
      }
    ]
  }
}
```

### VSCode Copilot settings

Open VSCode settings and add the following configuration to enable the MCP server for Copilot:
  ```json
  {
    // ...
    "mcp": {
      "servers": {
        "calculate": {
          "url": "http://localhost:8080/mcp",
        }
      }
    }
    // , ...
  }
  ```

### References
- [MCP Specs - Basic Lifecycle](https://modelcontextprotocol.io/specification/2025-03-26/basic/lifecycle)
- [MCP Specs - Tools](https://modelcontextprotocol.io/specification/2025-03-26/server/tools)
- [Visual testing tool for MCP](https://github.com/modelcontextprotocol/inspector)