package site

import (
	"os"
	"strconv"
	"strings"
)

type TreeEntry struct {
	Kind   string
	Name   string
	URL    string
	Target string
	Depth  int
}

func LoadTree(path string) ([]TreeEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []TreeEntry
	inTree := false
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "site_tree = [" {
			inTree = true
			continue
		}
		if !inTree {
			continue
		}
		if line == "]" {
			break
		}
		entry, ok := parseInlineTable(line)
		if ok {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func parseInlineTable(line string) (TreeEntry, bool) {
	line = strings.TrimSuffix(line, ",")
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
		return TreeEntry{}, false
	}

	line = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "{"), "}"))
	fields := splitFields(line)
	entry := TreeEntry{}

	for _, field := range fields {
		key, value, ok := strings.Cut(field, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "kind":
			entry.Kind = cleanValue(value)
		case "name":
			entry.Name = cleanValue(value)
		case "url":
			entry.URL = cleanValue(value)
		case "target":
			entry.Target = cleanValue(value)
		case "depth":
			depth, err := strconv.Atoi(value)
			if err == nil {
				entry.Depth = depth
			}
		}
	}

	return entry, entry.Kind != "" && entry.Name != ""
}

func splitFields(line string) []string {
	var fields []string
	var current strings.Builder
	inString := false

	for _, r := range line {
		switch r {
		case '"':
			inString = !inString
			current.WriteRune(r)
		case ',':
			if inString {
				current.WriteRune(r)
				continue
			}
			fields = append(fields, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		fields = append(fields, strings.TrimSpace(current.String()))
	}

	return fields
}

func cleanValue(value string) string {
	return strings.Trim(strings.TrimSpace(value), `"'`)
}
