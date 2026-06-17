//go:build !ignore

package ui

import _ "embed"

//go:embed design_tokens.css
var designTokensCSS string
// StyleSheet returns a single <style> tag containing the embedded design‑tokens.
func StyleSheet() string {
    return "<style>" + designTokensCSS + "</style>"
}
