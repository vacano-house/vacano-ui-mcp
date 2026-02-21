package docs

import (
	"path/filepath"
	"strings"
)

func Parse(files map[string]string, categoryMap CategoryMap) []DocEntry {
	var entries []DocEntry

	for path, content := range files {
		entry := parseFile(path, content, categoryMap)
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	return entries
}

func parseFile(path, content string, categoryMap CategoryMap) *DocEntry {
	// Skip index files
	base := filepath.Base(path)
	if base == "index.md" {
		return nil
	}

	// Skip .vitepress directory
	if strings.Contains(path, ".vitepress") {
		return nil
	}

	dir := filepath.Dir(path)

	if strings.HasSuffix(dir, "guide") {
		return parseGuide(content)
	}

	if strings.Contains(dir, "components") {
		return parseComponent(path, content, categoryMap)
	}

	return nil
}

func parseComponent(path, content string, categoryMap CategoryMap) *DocEntry {
	slug := strings.TrimSuffix(filepath.Base(path), ".md")

	name := extractH1(content)
	if name == "" {
		name = slugToName(slug)
	}

	category := CategoryUtility
	if cat, ok := categoryMap[slug]; ok {
		category = cat
	}

	return &DocEntry{
		Name:        name,
		Category:    category,
		Description: extractDescription(content),
		Content:     strings.TrimSpace(content),
	}
}

func parseGuide(content string) *DocEntry {
	name := extractH1(content)
	if name == "" {
		return nil
	}

	return &DocEntry{
		Name:        name,
		Category:    CategoryGuide,
		Description: extractDescription(content),
		Content:     strings.TrimSpace(content),
	}
}

func extractH1(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "## ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func extractDescription(section string) string {
	lines := strings.Split(section, "\n")
	pastHeading := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			pastHeading = true
			continue
		}

		if pastHeading {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "```") {
				break
			}
			return trimmed
		}
	}

	return ""
}

func slugToName(slug string) string {
	parts := strings.Split(slug, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
