package glassnode

import "github.com/c9s/requestgen"

//go:generate requestgen -method GET -type FeesRequest -url "/v1/metrics/fees/:metric" -responseType Response
type FeesRequest struct {
	Client requestgen.AuthenticatedAPIClient

	Asset           string   `param:"a,required,query"`
	Since           int64    `param:"s,query"`
	Until           int64    `param:"u,query"`
	Interval        Interval `param:"i,query"`
	Format          Format   `param:"f,query"`
	TimestampFormat string   `param:"timestamp_format,query"`

	Metric string `param:"metric,slug"`
}