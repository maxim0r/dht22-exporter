package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initMetrics(sensor Sensor, path string) error {

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "dht22",
			Name:      "temperature_celsius",
			Help:      "Temperature in Celsius",
		},
		func() float64 {
			_, temperature, err := sensor.Values()
			if err != nil {
				logger.Error("reading temperature error: %s", err)
			}
			return float64(temperature)
		},
	)); err != nil {
		return fmt.Errorf("error registering temperature gaugefunc: %w", err)
	}

	logger.Info("GaugeFunc 'temperature_celsius', registered.")

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "dht22",
			Name:      "humidity_percent",
			Help:      "Humidity in percent",
		},
		func() float64 {
			humidity, _, err := sensor.Values()
			if err != nil {
				logger.Error("error reading humidity %v", err)
			}
			return float64(humidity)
		},
	)); err != nil {
		return fmt.Errorf("error registering humidity gaugefunc: %w", err)
	}

	logger.Info("GaugeFunc 'temperature_celsius', registered.")

	http.Handle(path, promhttp.Handler())
}
