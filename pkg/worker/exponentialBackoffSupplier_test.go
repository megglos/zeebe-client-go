/*
 * Copyright Camunda Services GmbH and/or licensed to Camunda Services GmbH under
 * one or more contributor license agreements. See the NOTICE file distributed
 * with this work for additional information regarding copyright ownership.
 * Licensed under the Zeebe Community License 1.1. You may not use this file
 * except in compliance with the Zeebe Community License 1.1.
 */
package worker

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestExponentialBackoffSupplier_ShouldReturnDelayWithinBounds(t *testing.T) {
	iterations := 10

	minDelay := time.Millisecond * 50
	maxDelay := time.Second * 5
	e := NewExponentialBackoffBuilder().
		JitterFactor(0).
		MinDelay(minDelay).
		MaxDelay(maxDelay).
		BackoffFactor(1.6).
		Build()

	var retryDelays []time.Duration
	for i := 0; i < iterations; i++ {
		if len(retryDelays) == 0 {
			retryDelays = []time.Duration{e.SupplyRetryDelay(0)}
		}
		retryDelays = append(retryDelays, e.SupplyRetryDelay(retryDelays[i]))
	}

	// then - minDelay is equal to the first retryDelay
	// then - maxDelay is equal to the last retryDelay
	assert.Equal(t, minDelay, retryDelays[0])
	assert.Equal(t, maxDelay, retryDelays[10])
}

func TestExponentialBackoffSupplier_IsStrictlyIncreasing(t *testing.T) {
	iterations := 100
	t.Run("Backoff is strictly increasing", func(t *testing.T) {
		e := NewExponentialBackoffBuilder().
			JitterFactor(0).
			Build()
		var retryDelays []time.Duration
		for i := 0; i < iterations; i++ {
			if len(retryDelays) == 0 {
				retryDelays = []time.Duration{e.SupplyRetryDelay(0)}
			}
			retryDelays = append(retryDelays, e.SupplyRetryDelay(retryDelays[i]))
		}
		for i, delay := range retryDelays {
			// Skip first delay
			if i == 0 {
				continue
			}
			// then - as we used 0 for jitter factor, we can guarantee all are increasing or at least equal
			assert.GreaterOrEqual(t, delay, retryDelays[i-1], "backoff is strictly increasing")
		}

	})
}

func TestExponentialBackoffSupplier_ShouldBeRandomizedWithJitter(t *testing.T) {
	t.Run("backoff should be randomized with jitter", func(t *testing.T) {
		iterations := 100
		maxDelay := time.Second * 5
		minDelay := time.Millisecond * 50
		jitterFactor := 0.2
		e := NewExponentialBackoffBuilder().
			MaxDelay(maxDelay).
			MinDelay(minDelay).
			JitterFactor(jitterFactor).
			BackoffFactor(1.5).
			Build()

		maxDelayMillis := float64(maxDelay.Milliseconds())
		lowerMaxBound := math.Round(maxDelayMillis + maxDelayMillis*-jitterFactor)
		upperMaxBound := math.Round(maxDelayMillis + maxDelayMillis*jitterFactor)

		lowerMaxBoundDuration := time.Duration(lowerMaxBound * float64(time.Millisecond))
		upperMaxBoundDuration := time.Duration(upperMaxBound * float64(time.Millisecond))

		var retryDelays []time.Duration
		// when
		for i := 0; i < iterations; i++ {
			if len(retryDelays) == 0 {
				retryDelays = []time.Duration{e.SupplyRetryDelay(maxDelay)}
			}
			retryDelays = append(retryDelays, e.SupplyRetryDelay(retryDelays[i]))
		}

		// then
		for i, delay := range retryDelays {
			// retryDelay is in bounds
			betweenBounds := delay > lowerMaxBoundDuration && delay < upperMaxBoundDuration
			assert.True(t, betweenBounds, "is between lower and upper bound")

			// Skip first delay
			if i == 0 {
				continue
			}
			// then - as we used 0 for jitter factor, we can guarantee all are sorted
			assert.IsIncreasing(t, delay, retryDelays[i-1], "backoff is strictly increasing")
		}

	})
}
