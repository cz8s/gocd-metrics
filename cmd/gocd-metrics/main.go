package main

import (
	"flag"
	"log"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/cz8s/gocd-metrics/pkg/observe"
	"github.com/cz8s/gocd-metrics/pkg/web"
	"github.com/facebookgo/grace/gracehttp"
)

var (
	addr              = flag.String("addr", "0.0.0.0:9090", "Primary HTTP addr")
	debug             = flag.Bool("debug", false, "Set the logging level to debug")
	gocdHost          = flag.String("gocd-host", env("GOCD_HOST"), "The host for the GoCD server, including scheme and port.")
	gocdUsername      = flag.String("gocd-username", env("GOCD_USERNAME"), "The username for the GoCD user.")
	gocdPassword      = flag.String("gocd-password", env("GOCD_PASSWORD"), "The password for the GoCD user.")
	gocdSkipTLSVerify = flag.Bool("gocd-skip-tls-verify", envBool("GOCD_SKIP_TLS_VERIFY"), "Set to true to not verify the authenticity the GoCD server's TLS certificate.")

	writeTimeout = time.Second * 30
	readTimeout  = time.Second * 30
)

func main() {
	flag.Parse()

	r := web.NewRouter()
	r = observe.RegisterPrometheus(r)
	srv := &http.Server{
		Handler:      r,
		Addr:         *addr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
	gocdClient := observe.NewGocdClient(*gocdHost, *gocdUsername, *gocdPassword, *gocdSkipTLSVerify)
	metrics := observe.NewGocdMetrics()

	go func() {
		for {
			if err := observe.UpdateGocdMetrics(&metrics, gocdClient); err != nil {
				log.Print(err)
			}
			observe.UpdatePrometheus(metrics)
			time.Sleep(time.Duration(10 * time.Second))
		}
	}()

	gracehttp.Serve(srv)
}

func env(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return ""
	}
	return value
}

func envBool(name string) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		return false
	}
	return value == "true"
}
