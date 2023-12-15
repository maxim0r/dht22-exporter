package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/morus12/dht22"
)

var (
	listen      = flag.String("listen-address", ":9543", "The address to listen on for HTTP requests.")
	metricsPath = flag.String("metrics-path", "/metrics", "The path of the metrics endpoint.")
	gpioPort    = flag.String("gpio-port", "4", "The GPIO port where DHT22 is connected.")

	logger *slog.Logger
)

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	flag.Parse()

	sensor, err := initSensor(*gpioPort)
	if err != nil {
		logger.Error("init DHT22 sensor error: %s", err)
		return
	}

	if err := initMetrics(sensor, *metricsPath); err != nil {
		logger.Error("init metrics error: %s", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>DHT22 Exporter</title></head>
             <body>
             <h1>DHT22 Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	if err := http.ListenAndServe(*listen, nil); err != nil {
		logger.Error("listener error: %s", err)
	}
}
