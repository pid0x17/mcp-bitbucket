package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/pid0x17/mcp-bitbucket/internal/bitbucket"
)

type MCPServer struct {
	mcpServer *server.MCPServer
	bbClient  bitbucket.RepositoryFetcher
}

func NewMCPServer(bbClient bitbucket.RepositoryFetcher) *MCPServer {
	s := server.NewMCPServer(
		"Bitbucket MCP Server",
		"1.0.0",
	)

	mcpSrv := &MCPServer{
		mcpServer: s,
		bbClient:  bbClient,
	}

	mcpSrv.registerGetRepositoryTool()

	return mcpSrv
}

func (s *MCPServer) Serve() error {
	return server.ServeStdio(s.mcpServer)
}

func (s *MCPServer) registerGetRepositoryTool() {
	tool := mcp.NewTool("get_repository",
		mcp.WithDescription("Fetch details of a Bitbucket repository (name, description, privacy status)."),
		mcp.WithString("workspace", mcp.Required(), mcp.Description("The Bitbucket workspace ID")),
		mcp.WithString("repoSlug", mcp.Required(), mcp.Description("The repository slug/name")),
	)

	s.mcpServer.AddTool(tool, s.handleGetRepository)
}

func (s *MCPServer) handleGetRepository(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("arguments must be a JSON object"), nil
	}

	workspace, ok := args["workspace"].(string)
	if !ok {
		return mcp.NewToolResultError("missing or invalid 'workspace' argument"), nil
	}

	repoSlug, ok := args["repoSlug"].(string)
	if !ok {
		return mcp.NewToolResultError("missing or invalid 'repoSlug' argument"), nil
	}

	repo, err := s.bbClient.GetRepository(ctx, workspace, repoSlug)

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Bitbucket API error: %v", err)), nil
	}

	responseText := fmt.Sprintf("Repository: %s\nDescription: %s\nPrivate: %t",
		repo.Name, repo.Description, repo.IsPrivate)

	return mcp.NewToolResultText(responseText), nil
}
