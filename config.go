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

//go:generate paramgen -output=paramgen_src.go SourceConfig

package weather

import "time"

type SourceConfig struct {
	// url that contains the weather data
	URL string `json:"url" default:"https://api.openweathermap.org/data/2.5/weather"`
	// how often the connector will get data from the url
	PollingPeriod time.Duration `json:"pollingPeriod" default:"5m"`
	// city name to get the current weather for, ex: California, San Francisco, london. you can find the cities list
	// {city.list.json.gz} on http://bulk.openweathermap.org/sample/
	City string `json:"city" default:"new york"`
	// your unique API key (you can always find it on your account page under https://home.openweathermap.org/api_keys)
	APPID string `json:"appid" validate:"required"`
	// units of measurement, for Fahrenheit use imperial, for Celsius use metric, for Kelvin use standard.
	Units string `json:"units" default:"imperial" validate:"inclusion=imperial|standard|metric"`
}
