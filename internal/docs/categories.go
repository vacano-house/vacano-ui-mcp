package docs

import (
	"regexp"
	"strings"
)

type CategoryMap map[string]Category

var sectionRegex = regexp.MustCompile(`text:\s*'([^']+)'`)
var linkRegex = regexp.MustCompile(`link:\s*'(/(components|lib)/[^']+)'`)

func ParseCategories(configContent string) CategoryMap {
	cm := make(CategoryMap)

	if configContent == "" {
		return cm
	}

	sections := splitSidebarSections(configContent)

	for _, section := range sections {
		categoryName := extractSectionName(section)
		if categoryName == "" {
			continue
		}

		category := normalizeCategoryName(categoryName)
		slugs := extractComponentSlugs(section)

		for _, slug := range slugs {
			cm[slug] = category
		}
	}

	return cm
}

func splitSidebarSections(content string) []string {
	var allSections []string

	for _, key := range []string{"'/components/':", "'/lib/':"} {
		sections := extractSectionsForKey(content, key)
		allSections = append(allSections, sections...)
	}

	return allSections
}

func extractSectionsForKey(content, key string) []string {
	start := strings.Index(content, key)
	if start == -1 {
		return nil
	}

	// Find the opening bracket of the array
	arrayStart := strings.Index(content[start:], "[")
	if arrayStart == -1 {
		return nil
	}
	start += arrayStart

	// Find matching closing bracket
	depth := 0
	end := start
	for i := start; i < len(content); i++ {
		switch content[i] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				end = i + 1
				goto found
			}
		}
	}
	return nil

found:
	block := content[start:end]

	// Split by top-level objects in the array
	var sections []string
	depth = 0
	sectionStart := -1

	for i := 0; i < len(block); i++ {
		switch block[i] {
		case '{':
			if depth == 0 {
				sectionStart = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && sectionStart >= 0 {
				sections = append(sections, block[sectionStart:i+1])
				sectionStart = -1
			}
		}
	}

	return sections
}

func extractSectionName(section string) string {
	matches := sectionRegex.FindStringSubmatch(section)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func extractComponentSlugs(section string) []string {
	matches := linkRegex.FindAllStringSubmatch(section, -1)
	var slugs []string

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		// Extract slug from "/components/slug" or "/lib/slug"
		slug := match[1]
		slug = strings.TrimPrefix(slug, "/components/")
		slug = strings.TrimPrefix(slug, "/lib/")
		if slug != "" {
			slugs = append(slugs, slug)
		}
	}

	return slugs
}

func normalizeCategoryName(name string) Category {
	lower := strings.ToLower(name)
	normalized := strings.ReplaceAll(lower, " ", "-")

	switch normalized {
	case "form":
		return CategoryForm
	case "data-display":
		return CategoryDataDisplay
	case "feedback":
		return CategoryFeedback
	case "layout":
		return CategoryLayout
	case "navigation":
		return CategoryNavigation
	case "utility":
		return CategoryUtility
	case "overview":
		return CategoryOverview
	case "lib", "utils", "utilities", "hooks", "types", "constants":
		return CategoryLib
	default:
		return CategoryUtility
	}
}
