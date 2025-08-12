package checks

import (
	"strings"

	"github.com/gochelias/webaudit/internal/models"
	"github.com/gochelias/webaudit/internal/utils"
	"github.com/gocolly/colly"
)

func CheckSEO(c *colly.Collector, data *models.Report) {
	c.OnHTML("title", func(h *colly.HTMLElement) {
		if h.Text == "" {
			data.SEO = append(data.SEO, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: "The `title` tag has no content",
			})
		}
	})

	c.OnHTML("meta[name='description']", func(h *colly.HTMLElement) {
		if h.Attr("content") == "" {
			data.SEO = append(data.SEO, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: "The `description` meta tag has no content",
			})
		}
	})

	c.OnHTML("meta[name='robots']", func(h *colly.HTMLElement) {
		robots := h.Attr("content")

		isPagination := strings.Contains(utils.FormatURL(h.Request.URL.String()), "/page/")
		isCategory := strings.Contains(utils.FormatURL(h.Request.URL.String()), "/category/")

		if !strings.Contains(robots, "index, follow") && !isPagination && !isCategory {
			data.SEO = append(data.SEO, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: "Not indexing correctly",
			})
		} else if strings.Contains(robots, "index, follow") && isPagination && isCategory {
			data.SEO = append(data.SEO, models.Issue{
				Path:    utils.FormatURL(h.Request.URL.String()),
				Message: "It should not be indexed",
			})
		}
	})
}
