package utils

import (
	"veltiq/internal/core/domain"
)

func Start() {
	stats := make(chan domain.Metrics)
	go func() {
		for {
			select {
			case <-stats:
				processStat(stats)
			}
		}
	}()
}

func processStat(stats chan domain.Metrics) {
	
}