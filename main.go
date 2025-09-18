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

	// Register the filesystem_read_file tool
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_read_file",
		Description: "Reads the content of a file.",
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

	// Register the filesystem_write_file tool
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_write_file",
		Description: "Writes content to a file.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.WriteFileArgs) (*mcp.CallToolResult, any, error) {
		_, err := filesystem.WriteFile(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "File written successfully."},
			},
		}, nil, nil
	})

	// Register the filesystem_list_files tool
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_list_files",
		Description: "Lists files and directories in a path.",
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

	// Register the filesystem_file_exists tool
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "filesystem_file_exists",
		Description: "Checks if a file or directory exists.",
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

	// Register HTTP handler
	http.Handle("/mcp/", mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return srv
	}, nil))

	log.Println("MCPilot pair server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}