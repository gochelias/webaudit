package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gochelias/webaudit/internal/checks"
	"github.com/gochelias/webaudit/internal/models"
	"github.com/gochelias/webaudit/internal/report"
	"github.com/gochelias/webaudit/internal/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"

	"github.com/gocolly/colly"
	"github.com/slok/gospinner"
)

type Crawler struct {
	Data       *models.Report
	TotalPages int
	Spinner    *gospinner.Spinner
	baseURL    string
	collector  *colly.Collector
	config     models.Config
}

func Start(baseURL string, s *gospinner.Spinner, config models.Config) (*Crawler, error) {
	startTime := time.Now()
	s.SetMessage("Ready")

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := colly.NewCollector(
		colly.AllowedDomains(u.Host),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: config.Parallelism,
		Delay:       config.Delay,
	})

	s.Succeed()

	data := report.New()

	crawler := &Crawler{
		Data:      &data,
		Spinner:   s,
		baseURL:   baseURL,
		collector: c,
		config:    config,
	}

	s.Start("Starting")
	crawler.setup()

	c.Visit(u.String())

	elapsed := time.Since(startTime)
	s.SetMessage(fmt.Sprintf("%d pages scanned in %g", crawler.TotalPages, elapsed.Seconds()))
	s.Succeed()
	s.Finish()

	header := []string{"Status", "Issues with", "Total"}

	t := [][]string{
		{utils.GetStatus(crawler.Data.HTTPErrors, "🚨"), "HTTP", fmt.Sprintf("%d", len(crawler.Data.HTTPErrors))},
		{utils.GetStatus(crawler.Data.BrokenLinks, "❌"), "Links", fmt.Sprintf("%d", len(crawler.Data.BrokenLinks))},
		{utils.GetStatus(crawler.Data.HTMLStructure, "⚠️"), "HTML", fmt.Sprintf("%d", len(crawler.Data.HTMLStructure))},
		{utils.GetStatus(crawler.Data.SEO, "ℹ️"), "SEO", fmt.Sprintf("%d", len(crawler.Data.SEO))},
		{utils.GetStatus(crawler.Data.Accessibility, "ℹ️"), "Accessibility", fmt.Sprintf("%d", len(crawler.Data.Accessibility))},
		{utils.GetStatus(crawler.Data.Debug, "🔍"), "Debug", fmt.Sprintf("%d", len(crawler.Data.Debug))},
	}

	fmt.Println("")
	table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewMarkdown()))
	table.Header(header)
	table.Bulk(t)
	table.Render()

	return crawler, err
}

func (cr *Crawler) setup() {
	cr.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("x-vercel-protection-bypass", cr.config.VercelBypass)
		r.Headers.Set("x-vercel-set-bypass-cookie", "true")

		cr.TotalPages++
		cr.Spinner.SetMessage(fmt.Sprintf("(%d pages) %s", cr.TotalPages, utils.FormatURL(r.URL.String())))
	})

	cr.collector.RedirectHandler = func(req *http.Request, via []*http.Request) error {
		if req.URL.Path == "/" {
			return nil
		}

		cr.Data.Redirects = append(cr.Data.Redirects, models.Redirects{
			StatusCode: req.Response.StatusCode,
			Referer:    utils.FormatURL(req.Referer()),
			Path:       utils.FormatURL(req.URL.String()),
		})

		return nil
	}

	cr.collector.OnError(func(r *colly.Response, err error) {
		cr.Data.HTTPErrors = append(cr.Data.HTTPErrors, models.Issue{
			Path:    utils.FormatURL(r.Request.URL.String()),
			Message: fmt.Sprintf("%d %s", r.StatusCode, err.Error()),
		})
	})

	checks.CheckLinks(cr.collector, cr.Data)
	checks.CheckSEO(cr.collector, cr.Data)
	checks.CheckHTMLStructure(cr.collector, cr.Data)
	checks.CheckAccessibility(cr.collector, cr.Data)
}
