package checks

import (
	"fmt"
	"regexp"

	"github.com/gochelias/webaudit/internal/models"
	"github.com/gochelias/webaudit/internal/utils"
	"github.com/gocolly/colly"
)

func CheckHTMLStructure(c *colly.Collector, data *models.Report) {
	c.OnResponse(func(r *colly.Response) {
		tags := []string{"html", "head", "body", "main", "title"}

		for _, tag := range tags {
			rx := regexp.MustCompile(`(?i)<` + tag + `(?:\s|>)`)
			count := len(rx.FindAllString(string(r.Body), -1))

			if count > 1 {
				data.HTMLStructure = append(data.HTMLStructure, models.Issue{
					Path:    utils.FormatURL(r.Request.URL.String()),
					Message: fmt.Sprintf("%d `%s` tags found", count, tag),
				})
			} else if count == 0 {
				data.HTMLStructure = append(data.HTMLStructure, models.Issue{
					Path:    utils.FormatURL(r.Request.URL.String()),
					Message: fmt.Sprintf("No `%s` tag found", tag),
				})
			}
		}
	})
}
