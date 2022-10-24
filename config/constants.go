package config

import "time"

const (
	DatabaseQueryTimeLayout = `'YYYY-MM-DD"T"HH24:MI:SS"."MS"Z"TZ'`
	// DatabaseTimeLayout
	DatabaseTimeLayout string = time.RFC3339

	// Consumers
	ConsumerV1PhoneConsumer = "v1.phone.consumer"

	ErrorConsumer      = ".consumer1"
	InfoConsumer       = ".consumer2"
	DebugConsumer      = ".consumer3"
	AllMessageConsumer = ".consumer4"
)
