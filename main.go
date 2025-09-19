package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/seb-schulz/mcpilot-pair/middleware/auth"
	"github.com/seb-schulz/mcpilot-pair/tools/filesystem"
	"github.com/seb-schulz/mcpilot-pair/tools/make"
)

func prompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Hi prompt",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: "Say hi to " + req.Params.Arguments["name"]},
			},
		},
	}, nil
}

var embeddedResources = map[string]string{
	"info": "This is the hello example server.",
}

func embeddedResource(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	u, err := url.Parse(req.Params.URI)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "embedded" {
		return nil, fmt.Errorf("wrong scheme: %q", u.Scheme)
	}
	key := u.Opaque
	text, ok := embeddedResources[key]
	if !ok {
		return nil, fmt.Errorf("no embedded resource named %q", key)
	}
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{URI: req.Params.URI, MIMEType: "text/plain", Text: text},
		},
	}, nil
}

func main() {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "mcpilot-pair",
		Version: "0.3.0",
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

	// Registriere die Search-Funktion als Tool
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "search",
		Description: "Search for a regex pattern in files within the working directory. Returns a list of files with line numbers and matching lines.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.SearchArgs) (*mcp.CallToolResult, filesystem.SearchResult, error) {
		result, err := filesystem.Search(ctx, args)
		if err != nil {
			return nil, filesystem.SearchResult{}, fmt.Errorf("search failed: %v", err)
		}
		return &mcp.CallToolResult{}, result, nil
	})

	// Registriere fetch als Alias für filesystem_read_file
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "fetch",
		Description: "Alias for filesystem_read_file. Reads the content of a file within the working directory.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args filesystem.ReadFileArgs) (*mcp.CallToolResult, filesystem.ReadFileResult, error) {
		result, err := filesystem.ReadFile(ctx, args)
		if err != nil {
			return nil, filesystem.ReadFileResult{}, fmt.Errorf("fetch failed: %v", err)
		}
		return &mcp.CallToolResult{}, result, nil
	})

	srv.AddResource(&mcp.Resource{
		Name:     "info",
		MIMEType: "text/plain",
		URI:      "embedded:info",
	}, embeddedResource)

	srv.AddPrompt(&mcp.Prompt{Name: "greet"}, prompt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// MCP-Handler registrieren
	r.With(auth.APIKeyMiddleware).Handle("/mcp/*", mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return srv
	}, nil))

	log.Println("MCPilot pair server is running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
