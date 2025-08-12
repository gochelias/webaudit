package utils

import (
	"regexp"
)

func FormatURL(u string) string {
	r := regexp.MustCompile(`https.*?\.vercel.app`)

	return r.ReplaceAllString(u, "")
}
