package trafficstats

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import "net/http"
import "testing"
import "time"

import "github.com/apache/trafficcontrol/lib/go-tc"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

func TestConfigFromRequest(t *testing.T) {
	start := "2019-09-30T00:00:00Z"
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		t.Fatalf("Failed to parse test start time: %v", err)
	}
	end := "2019-10-01 00:00:00-07"
	endTime, err := time.Parse(tc.TimeLayout, end)
	if err != nil {
		t.Fatalf("Failed to parse test end time: %v", err)
	}

	inf := api.APIInfo{
		Params: map[string]string{
			"limit":           "10",
			"offset":          "0",
			"orderby":         "time",
			"startDate":       start,
			"endDate":         end,
			"interval":        "1m",
			"metricType":      "tps_total",
			"deliveryService": "test",
		},
	}

	r, e := http.NewRequest(http.MethodGet, "https://example.test/api/1.4/deliveryservice_stats", nil)
	if e != nil {
		t.Fatalf("Failed to build test request: %v", e)
	}
	r.Header.Add(tc.ContentType, tc.ApplicationJson)

	cfg, err, code := ConfigFromRequest(r, &inf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if code != http.StatusOK {
		t.Errorf("Expected OK status, but was %d", code)
	}
	if cfg.DeliveryService != "test" {
		t.Errorf("Expected config to have DS 'test', but was '%s'", cfg.DeliveryService)
	}
	if !cfg.End.Equal(endTime) {
		t.Errorf("Expected end time to be %v, but was %v", endTime, cfg.End)
	}
	if !cfg.Start.Equal(startTime) {
		t.Errorf("Expected start time to be %v, but was %v", startTime, cfg.Start)
	}
	if cfg.ExcludeSeries {
		t.Errorf("Expected series to not be excluded, but it was")
	}
	if cfg.ExcludeSummary {
		t.Errorf("Expected summary to not be excluded, but it was")
	}
	if cfg.Interval != tc.OneMinute {
		t.Errorf("Expected interval to be '1m', but it was %s", cfg.Interval)
	}
	if cfg.Limit == nil {
		t.Errorf("Expected limit to not be nil, but it was")
	} else if *cfg.Limit != 10 {
		t.Errorf("Expected limit to be 10, but it was %d", *cfg.Limit)
	}
	if cfg.MetricType != "tps_total" {
		t.Errorf("Expected metric type to be tps_total, but it was %s", cfg.MetricType)
	}
	if cfg.Offset == nil {
		t.Errorf("Expected offset to not be nil, but it was")
	} else if *cfg.Offset != 0 {
		t.Errorf("Expected offset to be 0, but it was %d", *cfg.Offset)
	}
	if cfg.OrderBy == nil {
		t.Errorf("Expected Order By to not be nil, but it was")
	} else if *cfg.OrderBy != tc.TimeOrder {
		t.Errorf("Expected Order by to be time, but it was %s", *cfg.OrderBy)
	}
	if cfg.Unix {
		t.Errorf("Expected Unix to not be set without MIME parameter, but it was")
	}
}
