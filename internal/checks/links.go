package checks

import (
	"fmt"
	"strings"

	"github.com/gochelias/webaudit/internal/models"
	"github.com/gochelias/webaudit/internal/utils"

	"github.com/gocolly/colly"
)

func CheckLinks(c *colly.Collector, data *models.Report) {
	var brokenLinks []models.Issue

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		href := h.Attr("href")

		if href == "/" || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "javascript:") || strings.Contains(href, "?") || strings.Contains(href, "#") {
			return
		}

		if href == "" || href == "#" || strings.HasPrefix(href, "-") {
			data.BrokenLinks = append(brokenLinks, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: fmt.Sprintf("`%s` invalid link", href),
			})
		}

		h.Request.Visit(href)
	})
}
