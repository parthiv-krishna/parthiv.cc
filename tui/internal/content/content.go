package content

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Page struct {
	Title       string
	Description string
	Path        string
	Section     string
	Body        string
}

func Load(root string) ([]Page, error) {
	var pages []Page

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		meta, body := splitFrontMatter(string(raw))
		if strings.EqualFold(meta["draft"], "true") {
			return nil
		}
		if strings.TrimSpace(body) == "" {
			return nil
		}

		section := filepath.Dir(rel)
		if section == "." {
			section = ""
		}

		page := Page{
			Title:       firstNonEmpty(meta["title"], titleFromBody(body), titleFromPath(rel)),
			Description: meta["description"],
			Path:        rel,
			Section:     section,
			Body:        strings.TrimSpace(body),
		}
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(pages) == 0 {
		return nil, errors.New("no markdown pages found")
	}

	sort.Slice(pages, func(i, j int) bool {
		if pages[i].Path == "_index.md" {
			return true
		}
		if pages[j].Path == "_index.md" {
			return false
		}
		return pages[i].Title < pages[j].Title
	})

	return pages, nil
}

func splitFrontMatter(src string) (map[string]string, string) {
	meta := map[string]string{}
	lines := strings.Split(src, "\n")
	if len(lines) == 0 {
		return meta, src
	}

	delim := strings.TrimSpace(lines[0])
	if delim != "+++" && delim != "---" {
		return meta, src
	}

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == delim {
			return meta, strings.Join(lines[i+1:], "\n")
		}
		key, value, ok := strings.Cut(lines[i], "=")
		if !ok {
			key, value, ok = strings.Cut(lines[i], ":")
		}
		if !ok {
			continue
		}
		meta[strings.TrimSpace(key)] = cleanValue(value)
	}

	return meta, src
}

func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return value
}

func titleFromBody(body string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func titleFromPath(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	return strings.TrimSpace(base)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return "Untitled"
}
