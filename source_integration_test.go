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
	"errors"
	"os"
	"testing"
	"time"

	weather "github.com/conduitio-labs/conduit-connector-weather"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

func TestSource_GetWeather(t *testing.T) {
	APPID := os.Getenv("APPID")
	if APPID == "" {
		t.Skipf("APPID env var must be set")
	}

	is := is.New(t)
	con := weather.Source{}
	err := con.Configure(context.Background(), map[string]string{
		"appid":         APPID,
		"city":          "london",
		"units":         "metric",
		"url":           "https://api.openweathermap.org/data/2.5/weather",
		"pollingPeriod": "3s",
	})
	is.NoErr(err)
	ctx := context.Background()
	err = con.Open(ctx, sdk.Position{})
	is.NoErr(err)
	// first read should succeed
	_, err = con.Read(ctx)
	is.NoErr(err)
	// it hasn't been 3 seconds yet, second read should fail
	rec, err := con.Read(ctx)
	is.True(errors.Is(err, sdk.ErrBackoffRetry))
	is.Equal(rec, sdk.Record{})
	// delay 3 seconds, read should work now
	time.Sleep(3 * time.Second)
	_, err = con.Read(ctx)
	is.NoErr(err)
}
