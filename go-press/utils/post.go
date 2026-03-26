package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(title string) string {
	// lowercase
	slug := strings.ToLower(title)

	// replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// remove special characters (keep a-z, 0-9, "-")
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// remove multiple dash
	regDash := regexp.MustCompile(`-+`)
	slug = regDash.ReplaceAllString(slug, "-")

	// trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}