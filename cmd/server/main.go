package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vacano-house/vacano-ui-mcp/internal/config"
	"github.com/vacano-house/vacano-ui-mcp/internal/docs"
	"github.com/vacano-house/vacano-ui-mcp/internal/repo"
	"github.com/vacano-house/vacano-ui-mcp/internal/tools"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Git repo
	repository, err := repo.New(cfg.Repo)
	if err != nil {
		log.Fatalf("Failed to init repo: %v", err)
	}
	defer repository.Cleanup()

	// Clone repo
	log.Println("Cloning repository...")
	if err := repository.Clone(); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	// Docs store
	store := docs.NewStore()

	// Initial docs parse
	if err := refreshDocs(repository, store); err != nil {
		log.Fatalf("Failed to parse documentation: %v", err)
	}
	log.Println("Documentation loaded successfully")

	// Background refresh
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startRefreshLoop(ctx, cfg.Docs.RefreshInterval, repository, store)

	// MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "vacano-ui-docs",
			Version: "1.0.0",
		},
		nil,
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_docs",
		Description: "Search across all vacano-ui documentation by keyword. Searches in component names, descriptions, and full content.",
	}, tools.NewSearchHandler(store))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_component_docs",
		Description: "Get full documentation for a specific vacano-ui component by exact name (e.g. Button, Modal, DatePicker).",
	}, tools.NewGetComponentHandler(store))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_components",
		Description: "List all available vacano-ui components. Optionally filter by category: form, data-display, feedback, layout, navigation, utility, overview, guide.",
	}, tools.NewListHandler(store))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_icons",
		Description: "Search vacano-ui icons (1,894 Lucide icons) by name, description, or category. Icons are imported from '@vacano/ui/icons'. Use this to find the right icon for a UI element.",
	}, tools.NewSearchIconsHandler(store))

	// Streamable HTTP handler
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	// HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: mux,
	}

	go func() {
		log.Printf("MCP server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func refreshDocs(repository *repo.Repo, store *docs.Store) error {
	// Fetch VitePress config for category mapping
	configContent, err := repository.FetchVitePressConfig()
	if err != nil {
		log.Printf("Warning: failed to read VitePress config: %v", err)
	}
	categoryMap := docs.ParseCategories(configContent)

	// Fetch and parse docs
	files, err := repository.FetchDocs()
	if err != nil {
		return fmt.Errorf("failed to read docs: %w", err)
	}

	entries := docs.Parse(files, categoryMap)
	store.Reload(entries)

	// Parse icons from icons.md
	for path, content := range files {
		if strings.HasSuffix(path, "components/icons.md") {
			icons := docs.ParseIcons(content)
			store.ReloadIcons(icons)
			log.Printf("Loaded %d icons", len(icons))
			break
		}
	}

	log.Printf("Loaded %d documentation entries", len(entries))
	return nil
}

func startRefreshLoop(ctx context.Context, interval time.Duration, repository *repo.Repo, store *docs.Store) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("Refreshing documentation...")
			if err := repository.Pull(); err != nil {
				log.Printf("Failed to pull: %v", err)
				continue
			}
			if err := refreshDocs(repository, store); err != nil {
				log.Printf("Failed to refresh docs: %v", err)
			}
		}
	}
}
