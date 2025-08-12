package checks

import (
	"fmt"
	"strings"

	"github.com/gochelias/webaudit/internal/models"
	"github.com/gochelias/webaudit/internal/utils"
	"github.com/gocolly/colly"
)

func CheckAccessibility(c *colly.Collector, data *models.Report) {
	c.OnHTML("img", func(h *colly.HTMLElement) {
		if h.Attr("alt") == "" {
			data.Accessibility = append(data.Accessibility, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: "Missing alternative text",
				Src:     h.Attr("src"),
			})
		}
	})

	// https://dequeuniversity.com/rules/axe/4.10/link-name
	c.OnHTML("a", func(h *colly.HTMLElement) {
		content := strings.Trim(h.Text, "\n ")
		label := strings.Trim(h.Attr("aria-label"), " ")

		if label == "" && content == "" {
			data.Accessibility = append(data.Accessibility, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: fmt.Sprintf("The link `%s` has no accessible content", h.Attr("href")),
			})
		}
	})
}
