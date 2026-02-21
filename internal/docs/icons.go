package docs

import (
	"strings"
)

func ParseIcons(content string) []IconEntry {
	if content == "" {
		return nil
	}

	var entries []IconEntry
	var currentCategory string

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track icon category headings (### Actions, ### Arrows and Chevrons, etc.)
		if strings.HasPrefix(trimmed, "### ") {
			cat := strings.TrimPrefix(trimmed, "### ")
			// Skip usage sections (Basic, Sizing, Coloring, etc.) — they come before "Icon Reference"
			switch cat {
			case "Basic", "Sizing", "Coloring", "With Buttons", "Inline with Text":
				currentCategory = ""
				continue
			}
			currentCategory = cat
			continue
		}

		// Parse table rows: | `IconName` | description |
		if currentCategory != "" && strings.HasPrefix(trimmed, "| `") {
			entry := parseIconRow(trimmed, currentCategory)
			if entry != nil {
				entries = append(entries, *entry)
			}
		}
	}

	return entries
}

func parseIconRow(line, category string) *IconEntry {
	// Format: | `IconName` | Description text |
	parts := strings.Split(line, "|")
	if len(parts) < 3 {
		return nil
	}

	namePart := strings.TrimSpace(parts[1])
	descPart := strings.TrimSpace(parts[2])

	// Extract name from backticks
	name := strings.Trim(namePart, "`")
	if name == "" || name == "Icon" {
		return nil
	}

	return &IconEntry{
		Name:        name,
		Description: descPart,
		Category:    category,
	}
}
