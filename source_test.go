// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package weather_test

import (
	"context"
	"testing"

	weather "github.com/conduitio-labs/conduit-connector-weather"
	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/matryer/is"
)

func TestSource_ConfigureAndOpen(t *testing.T) {
	is := is.New(t)
	con := weather.Source{}
	err := con.Configure(context.Background(), map[string]string{
		"appid": "my-id",
		"city":  "london",
		"units": "metric",
		"url":   "https://api.openweathermap.org/data/2.5/weather",
	})
	is.NoErr(err)
	ctx := context.Background()
	err = con.Open(ctx, opencdc.Position{})
	is.True(err != nil) // auth failed
	url := con.CreateRequestURL()
	is.Equal(url, "https://api.openweathermap.org/data/2.5/weather?q=london&APPID=my-id&units=metric")
}

func TestTeardownSource_NoOpen(t *testing.T) {
	con := weather.NewSource()
	err := con.Teardown(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
