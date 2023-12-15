package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/MichaelS11/go-dht"
)

type Sensor interface {
	Values() (float64, float64, error)
}

type dht22 struct {
	dht *dht.DHT

	temp, hum float64
	last time.Time
	sync.Mutex
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
	s.Lock()
	defer s.Unlock()
	
	if time.Since(s.last) < time.Second {
		return s.temp, s.hum, nil
	}
	
	t, h, err := s.dht.ReadRetry(10)
	if err != nil {
		return t, h, err
	}

	s.temp = t
	s.hum = h
	s.last = time.Now()
	 
	return s.temp, s.hum, nil
}
