package ui_test

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
)

// canonicalHTML parses an HTML fragment and returns a canonical string.
// It removes the automatic <html><head><body> wrapper that html.Parse adds.
func canonicalHTML(s string) (string, error) {
	ctx := &html.Node{Type: html.ElementNode, Data: "body"}
	nodes, err := html.ParseFragment(strings.NewReader(s), ctx)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, n := range nodes {
		if err := html.Render(&buf, n); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
