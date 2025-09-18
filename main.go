package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/seb-schulz/mcpilot-pair/filesystem"
	"github.com/seb-schulz/mcpilot-pair/middleware/auth"
	"github.com/seb-schulz/mcpilot-pair/tools/make"
)

func main() {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "mcpilot-pair",
		Version: "0.2.0",
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
		log.Println("filesystem_file_exists", result, err)
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

	// Register the make_run tool
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "make_run",
		Description: "Executes `make -C <directory> <target>` in a bash environment. Allowed targets: all, build, test, clean.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args make.RunMakeArgs) (*mcp.CallToolResult, any, error) {
		result, err := make.RunMake(ctx, args)
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// MCP-Handler registrieren
	r.With(auth.APIKeyMiddleware).Handle("/mcp/*", mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		log.Println(req.Header)
		return srv
	}, nil))

	log.Println("MCPilot pair server is running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
