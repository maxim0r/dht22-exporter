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

	dht, err := dht.NewDHT(p, dht.Celsius, "")
	if err != nil {
		return nil, fmt.Errorf("create DHT instance error: %w", err)
	}

	d := &dht22{
		dht: dht,
	}
	
	// data fetch worker
	go func() {
		t := time.NewTicker(time.Second*30)
		defer t.Stop()
		for range t.C {
			d.Lock()
			t, h, err := d.dht.ReadRetry(10)
			if err != nil {
			    d.Unlock()
				continue
			}
			d.temp = t
			d.hum = h
			d.last = time.Now()
			d.Unlock()
		}
	}()
	
	return d, nil
}

func (d *dht22) Values() (float64, float64, error) {
	d.Lock()
	defer d.Unlock()
	
	if time.Since(d.last) < time.Minute {
		return d.temp, d.hum, nil
	}
	
	t, h, err := d.dht.ReadRetry(10)
	if err != nil {
		return t, h, err
	}

	d.temp = t
	d.hum = h
	d.last = time.Now()
	 
	return d.temp, d.hum, nil
}
