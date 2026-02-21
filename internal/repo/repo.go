package repo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vacano-house/vacano-ui-mcp/internal/config"
)

type Repo struct {
	url        string
	branch     string
	sshKeyFile string
	localPath  string
}

func New(cfg config.RepoConfig) (*Repo, error) {
	tmpDir, err := os.MkdirTemp("", "vacano-ui-docs-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	r := &Repo{
		url:       cfg.URL,
		branch:    cfg.Branch,
		localPath: tmpDir,
	}

	if cfg.SSHKey != "" {
		keyFile, err := writeSSHKey(cfg.SSHKey)
		if err != nil {
			return nil, fmt.Errorf("failed to write SSH key: %w", err)
		}
		r.sshKeyFile = keyFile
	}

	return r, nil
}

func (r *Repo) Clone() error {
	args := []string{
		"clone",
		"--depth", "1",
		"--branch", r.branch,
		"--single-branch",
		r.url,
		r.localPath,
	}

	cmd := exec.Command("git", args...)
	r.setSSHEnv(cmd)

	log.Printf("Cloning %s (branch: %s) into %s", r.url, r.branch, r.localPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (r *Repo) Pull() error {
	cmd := exec.Command("git", "-C", r.localPath, "pull", "--ff-only")
	r.setSSHEnv(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %w, output: %s", err, string(output))
	}

	log.Printf("Git pull: %s", strings.TrimSpace(string(output)))
	return nil
}

func (r *Repo) FetchDocs() (map[string]string, error) {
	docsPath := filepath.Join(r.localPath, "docs")

	if _, err := os.Stat(docsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("docs directory not found at %s", docsPath)
	}

	files := make(map[string]string)

	err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Skip .vitepress directory
		relPath, _ := filepath.Rel(r.localPath, path)
		if strings.Contains(relPath, ".vitepress") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		files[relPath] = string(content)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk docs directory: %w", err)
	}

	return files, nil
}

func (r *Repo) FetchVitePressConfig() (string, error) {
	configPath := filepath.Join(r.localPath, "docs", ".vitepress", "config.ts")

	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("VitePress config not found, using default categories")
			return "", nil
		}
		return "", fmt.Errorf("failed to read VitePress config: %w", err)
	}

	return string(content), nil
}

func (r *Repo) Cleanup() {
	os.RemoveAll(r.localPath)
	if r.sshKeyFile != "" {
		os.Remove(r.sshKeyFile)
	}
}

func (r *Repo) setSSHEnv(cmd *exec.Cmd) {
	if r.sshKeyFile != "" {
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s -o StrictHostKeyChecking=no", r.sshKeyFile),
		)
	}
}

func writeSSHKey(key string) (string, error) {
	// godotenv may store literal \n instead of newlines
	normalized := strings.ReplaceAll(key, "\\n", "\n")

	// ensure trailing newline (required by SSH)
	if !strings.HasSuffix(normalized, "\n") {
		normalized += "\n"
	}

	f, err := os.CreateTemp("", "ssh-key-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := f.WriteString(normalized); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", fmt.Errorf("failed to write key: %w", err)
	}
	f.Close()

	if err := os.Chmod(f.Name(), 0600); err != nil {
		os.Remove(f.Name())
		return "", fmt.Errorf("failed to set key permissions: %w", err)
	}

	return f.Name(), nil
}
