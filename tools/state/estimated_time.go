package main

import "time"

const conversionFactorPercentage = 100

type estimatedTimeOfArrival struct {
	totalItems       int
	processedItems   int
	startTime        time.Time
	previousStepTime time.Time
}

func (e *estimatedTimeOfArrival) start(totalItems int) {
	e.totalItems = totalItems
	e.processedItems = 0
	e.startTime = time.Now()
	e.previousStepTime = e.startTime
}

// return eta time.Duration, percent float64, itemsPerSecond float64
func (e *estimatedTimeOfArrival) step(itemsProcessedInthisStep int) (time.Duration, float64, float64) {
	e.processedItems += itemsProcessedInthisStep

	curentTime := time.Now()
	elapsedTime := curentTime.Sub(e.startTime)
	eta := time.Duration(float64(elapsedTime) / float64(e.processedItems) * float64(e.totalItems-e.processedItems))
	percent := float64(e.processedItems) / float64(e.totalItems) * conversionFactorPercentage
	itemsPerSecond := float64(e.processedItems) / elapsedTime.Seconds()
	e.previousStepTime = curentTime
	return eta, percent, itemsPerSecond
}
