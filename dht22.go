package main

import (
	"fmt"

	"github.com/MichaelS11/go-dht"
)

type Sensor interface {
	Values() (float64, float64, error)
}

type dht22 struct {
	dht *dht.DHT
}

func initSensor(p string) (Sensor, error) {
	err := dht.HostInit()
	if err != nil {
		return nil, fmt.Errorf("HostInit error: %w", err)
	}

	d, err := dht.NewDHT(p, dht.Celsius, "")
	if err != nil {
		return nil, fmt.Errorf("create DHT instance error: %w", err)
	}

	return &dht22{
		dht: d,
	}, nil
}

func (s *dht22) Values() (float64, float64, error) {
	return s.dht.Read()
}
