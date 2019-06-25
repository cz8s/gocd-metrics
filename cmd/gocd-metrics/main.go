package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cz8s/gocd-metrics/pkg/observe"
	"github.com/cz8s/gocd-metrics/pkg/web"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/spf13/cobra"
)

var (
	metricsAddress            string
	gocdAddress               string
	gocdUsername              string
	gocdPassword              string
	gocdSkipTLSVerify         bool
	gocdScrapeIntervalSeconds = 10
	writeTimeout              = time.Second * 30
	readTimeout               = time.Second * 30
	rootCmd                   = &cobra.Command{
		Use:   "gocd-metrics",
		Short: "Exports GoCD build metrics to prometheus.",
		Long:  `Exports GoCD build metrics to prometheus.`,
		Run: func(cmd *cobra.Command, args []string) {
			r := web.NewRouter()
			r = observe.RegisterPrometheus(r)
			srv := &http.Server{
				Handler:      r,
				Addr:         metricsAddress,
				WriteTimeout: writeTimeout,
				ReadTimeout:  readTimeout,
			}
			gocdClient := observe.NewGocdClient(gocdAddress, gocdUsername, gocdPassword, gocdSkipTLSVerify)
			metrics := observe.NewGocdMetrics()

			go func() {
				for {
					if err := observe.UpdateGocdMetrics(&metrics, gocdClient); err != nil {
						log.Print(err)
					}
					observe.UpdatePrometheus(metrics)
					time.Sleep(time.Duration(time.Second * time.Duration(gocdScrapeIntervalSeconds)))
				}
			}()

			gracehttp.Serve(srv)
		},
	}
)

func main() {
	rootCmd.Flags().StringVar(&metricsAddress, "addr", "0.0.0.0:9090", "The address for exposing the metrics server")
	rootCmd.Flags().StringVar(&gocdAddress, "gocd-url", "", "The address for the GoCD server, including scheme and port")
	rootCmd.MarkFlagRequired("gocd-url")
	rootCmd.Flags().StringVar(&gocdUsername, "gocd-username", "", "The username for the GoCD user")
	rootCmd.Flags().StringVar(&gocdPassword, "gocd-password", "", "The password for the GoCD user")
	rootCmd.Flags().BoolVar(&gocdSkipTLSVerify, "gocd-skip-tls-verify", false, "Skip verification of the authenticity of the GoCD server's TLS certificate")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
