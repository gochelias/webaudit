package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/gochelias/webaudit/internal/crawler"
	"github.com/gochelias/webaudit/internal/models"
	"github.com/gochelias/webaudit/internal/report"
	"github.com/slok/gospinner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "webaudit",
	Short: "Audit and analyze websites for common issues.",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		bypass, _ := cmd.Flags().GetString("vercel-bypass")
		parallelism, _ := cmd.Flags().GetInt("parallelism")
		delay, _ := cmd.Flags().GetInt("delay")

		figure.NewFigure("webaudit", "smslant", true).Print()
		fmt.Println("")

		s, _ := gospinner.NewSpinner(gospinner.BouncingBar)
		s.Start("Configuration")

		cr, err := crawler.Start(url, s, models.Config{
			Parallelism:  parallelism,
			Delay:        time.Duration(delay) * time.Millisecond,
			VercelBypass: bypass,
		})
		if err != nil {
			log.Fatal(err)
		}

		report.Save("report.json", cr.Data)

		if len(cr.Data.HTTPErrors) > 0 || len(cr.Data.BrokenLinks) > 0 {
			fmt.Println("")
			log.Fatal("Problems were identified")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("url", "", "URL of the site to audit")
	rootCmd.Flags().String("vercel-bypass", "", "Secret to bypass Vercel Protection (VERCEL_AUTOMATION_BYPASS_SECRET)")
	rootCmd.Flags().Int("parallelism", 1, "Number of the maximum allowed concurrent requests")
	rootCmd.Flags().Int("delay", 1000, "Delay between requests in milliseconds")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("vercel-bypass")
}
