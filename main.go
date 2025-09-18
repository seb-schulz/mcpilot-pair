package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/seb-schulz/mcpilot-pair/filesystem"
)

func main() {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "mcpilot-pair",
		Version: "0.1.0",
	}, nil)

	// filesystem_read_file-Tool registrieren
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_read_file",
		Description: "Liest den Inhalt einer Datei.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.ReadFileArgs) (*mcp.CallToolResult, any, error) {
		result, err := filesystem.ReadFile(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: result.Content},
			},
		}, nil, nil
	})

	// filesystem_write_file-Tool registrieren
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_write_file",
		Description: "Schreibt Inhalt in eine Datei.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.WriteFileArgs) (*mcp.CallToolResult, any, error) {
		_, err := filesystem.WriteFile(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Datei erfolgreich geschrieben."},
			},
		}, nil, nil
	})

	// filesystem_list_files-Tool registrieren
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_list_files",
		Description: "Listet Dateien und Verzeichnisse in einem Pfad auf.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.ListFilesArgs) (*mcp.CallToolResult, any, error) {
		result, err := filesystem.ListFiles(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		jsonData, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonData)},
			},
		}, nil, nil
	})

	// filesystem_file_exists-Tool registrieren
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_file_exists",
		Description: "Prüft, ob eine Datei oder ein Verzeichnis existiert.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.FileExistsArgs) (*mcp.CallToolResult, any, error) {
		result, err := filesystem.FileExists(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		jsonData, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonData)},
			},
		}, nil, nil
	})

	// HTTP-Handler registrieren
	http.Handle("/mcp/", mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return srv
	}, nil))

	log.Println("MCPilot pair Server ist bereit auf :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Serverfehler: %v", err)
	}
}
