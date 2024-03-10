package helpers

import (
	"regexp"
)

func ExtractJavaScriptFromHTML(str string) string {
	regex, err := regexp.Compile(`(?is)\s*<script\b[^<]*>(.*?)<\/script>\s*`)
	if err != nil {
		println("Error extracting JS from HTML")
	}

	matches := regex.FindAllStringSubmatch(str, -1)
	if matches == nil {
		return ""
	}
	var scripts string
	for _, match := range matches {
		scripts += match[1]
	}

	return scripts
}
