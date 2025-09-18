# MCPilot Pair

**MCPilot Pair securely links your local dev environment to LLMs via reverse tunnels (SSH/Ngrok). Edit files, run Git commands, or execute `make test` - all while keeping control. Built on the Model Context Protocol (MCP), it’s your AI-ready co-pilot for private, efficient workflows.**

In case you want to know more about MCP, check out <https://modelcontextprotocol.io/docs/getting-started/intro>

## Installation

Install the tool using `go install`:

```bash
go install github.com/seb-schulz/mcpilot-pair@latest
```

Start the server:

```bash
mcpilot-pair
```

The API key is automatically saved in `.mcpilot-pair-api-key.txt`.

## Usage

### As a Server

```bash
mcpilot-pair  # Starts the MCP server on port 8080
```

### Establishing a Connection

You can connect to the MCP server in various ways:

#### Via SSH Tunnel

```bash
ssh -R 127.0.0.1:30204:127.0.0.1:8080 example.com
```

#### Using Apache Rewrite Rule

If you use Apache, you can use a Rewrite Rule:

```apache
RewriteCond %{HTTP_HOST} ^rand-sub\.example\.com$ [NC]
RewriteRule ^(.*) http://127.0.0.1:30204/$1 [proxy,last]
```

#### With Custom Connectors

The MCP server can be connected to various LLMs, including:

- **Le Chat (Mistral AI)** – supports custom MCP servers directly.
- **Claude (Anthropic)** – supports local and remote MCP servers.
- **Gemini (Google)** – supports MCP servers via Gemini CLI (don't know if this is useful).
- **ChatGPT:** – not supported yet.
  ChatGPT supports MCP servers only in a limited capacity. Check their documentation about details.

## Contributing

Pull requests are welcome! Make sure to:

1. Run `go test ./...` before committing.
2. Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/).

## License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.
