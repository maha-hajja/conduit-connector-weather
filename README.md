# Conduit Connector for Weather
[Conduit](https://conduit.io) source connector for weather.

## How to build?
Run `make` to build the connector.

## Testing
Run `make test` to run all the unit tests.

## Source
This source connector pulls data from [OpenWeather](https://openweathermap.org/) and pushes the weather data to 
downstream resources via Conduit.

### Configuration

| name            | description                           | required | default value |
|-----------------|---------------------------------------|----------|---------------|
| `appid`         | your unique API key (you can always find it on your account page under https://home.openweathermap.org/api_keys) | true     |           |
| `city`          | city name to get the current weather for, ex: California, San Francisco, london. you can find the cities list {city.list.json.gz} on http://bulk.openweathermap.org/sample/ | false | New York |
| `units`         | units of measurement, for Fahrenheit use imperial, for Celsius use metric, for Kelvin use standard.             | false     | imperial  |
| `url`           | url that contains the weather data                                                                              | false     | https://api.openweathermap.org/data/2.5/weather |
| `pollingPeriod` | how often the connector will get data from the url, formatted as a time.Duration string                         | false     | 5m        |

## Example
here's a pipeline configuration file sample:
```yaml
   pipelines:
   weather-pipeline:
     status: running
     name: weather-pipeline
     description: get the current weather in California every 5 minutes
     connectors:
       con-weather:
         type: source
         plugin: standalone:weather
         name: weather-source
         settings:
           city: California
           appid: ${APPID}
           pollingPeriod: 2m
           units: metric
       con-file:
         type: destination
         plugin: builtin:file
         name: file-dest
         settings:
           path: ./weather.txt
           sdk.record.format: template
           sdk.record.format.options: '{{ toJson .Payload.After }}'
```
make sure to export your appid key to the env variable $APPID before running conduit.
check [Pipeline Configuration Files Docs](https://github.com/ConduitIO/conduit/blob/main/docs/pipeline_configuration_files.md) to run this pipeline.

Results:
the output file `weather.txt` will have a new weather reading every two minute, and would look something like:
```json
{
  "base": "stations",
  "clouds": {
    "all": 100
  },
  "cod": 200,
  "coord": {
    "lat": 38.3004,
    "lon": -76.5074
  },
  "dt": 1675955436,
  "id": 4350049,
  "main": {
    "feels_like": 11.48,
    "humidity": 78,
    "pressure": 1022,
    "temp": 12.17,
    "temp_max": 13.94,
    "temp_min": 8.84
  },
  "name": "California",
  "sys": {
    "country": "US",
    "id": 2011802,
    "sunrise": 1675944220,
    "sunset": 1675982215,
    "type": 2
  },
  "timezone": -18000,
  "visibility": 10000,
  "weather": [
    {
      "description": "overcast clouds",
      "icon": "04d",
      "id": 804,
      "main": "Clouds"
    }
  ],
  "wind": {
    "deg": 150,
    "gust": 9.26,
    "speed": 6.17
  }
} 
```
