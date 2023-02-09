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

package weather

import (
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

var exampleConfig = map[string]string{
	"appid":         "my-appid-key",
	"city":          "San Francisco",
	"pollingPeriod": "10s",
	"units":         "metric",
}

func TestParseConfig(t *testing.T) {
	is := is.New(t)
	var got SourceConfig
	err := sdk.Util.ParseConfig(exampleConfig, &got)
	want := SourceConfig{
		APPID:         "my-appid-key",
		City:          "San Francisco",
		Units:         "metric",
		PollingPeriod: 10 * time.Second,
	}
	is.NoErr(err)
	is.Equal(want, got)
}
