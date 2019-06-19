package main

import (
	"flag"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/cz8s/gocd-metrics/pkg/observe"
	"github.com/cz8s/gocd-metrics/pkg/web"
	"github.com/facebookgo/grace/gracehttp"
)

var (
	addr  = flag.String("addr", "0.0.0.0:9090", "Primary HTTP addr")
	debug = flag.Bool("debug", false, "Set the logging level to debug")

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
	gracehttp.Serve(srv)
}
