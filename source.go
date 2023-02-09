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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"golang.org/x/time/rate"
)

type Source struct {
	sdk.UnimplementedSource

	config  SourceConfig
	client  *http.Client
	url     string
	limiter *rate.Limiter
}

func NewSource() sdk.Source {
	return &Source{}
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Weather Source Connector...")
	var config SourceConfig
	err := sdk.Util.ParseConfig(cfg, &config)
	if err != nil {
		return err
	}
	s.config = config
	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) error {
	s.client = &http.Client{}
	s.url = s.CreateRequestURL()
	// try pinging the URL with APPID
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, s.url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request %q: %w", s.config.URL, err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error pinging URL %q: %w", s.config.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authorization failed, check your APPID key")
	}

	s.limiter = rate.NewLimiter(rate.Every(s.config.PollingPeriod), 1)

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	err := s.limiter.Wait(ctx)
	if err != nil {
		return sdk.Record{}, err
	}
	rec, err := s.getRecord(ctx)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("error getting the weather data: %w", err)
	}
	return rec, nil
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	sdk.Logger(ctx).Debug().Str("position", string(position)).Msg("got ack")
	return nil // no ack needed
}

func (s *Source) Teardown(ctx context.Context) error {
	if s.client != nil {
		s.client.CloseIdleConnections()
	}
	return nil
}

func (s *Source) CreateRequestURL() string {
	return s.config.URL + "?" + "q=" + s.config.City + "&" + "APPID=" + s.config.APPID + "&" + "units=" + s.config.Units
}

func (s *Source) getRecord(ctx context.Context) (sdk.Record, error) {
	// create GET request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("error creating HTTP request: %w", err)
	}
	// get response
	resp, err := s.client.Do(req)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("error getting data from URL: %w", err)
	}
	defer resp.Body.Close()
	// check response status
	if resp.StatusCode != http.StatusOK {
		return sdk.Record{}, fmt.Errorf("response status should be %v, got status=%v", http.StatusOK, resp.StatusCode)
	}
	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("error reading body for response %v: %w", resp, err)
	}
	// parse json
	var structData sdk.StructuredData
	err = json.Unmarshal(body, &structData)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("failed to unmarshal body as JSON: %w", err)
	}
	// create record
	now := time.Now().Unix()
	timestamp, ok := structData["dt"]
	if !ok {
		return sdk.Record{}, fmt.Errorf("dt field not found in record: %w", err)
	}
	rec := sdk.Record{
		Payload: sdk.Change{
			Before: nil,
			After:  structData,
		},
		Operation: sdk.OperationCreate,
		Position:  sdk.Position(fmt.Sprintf("unix-%v", now)),
		Key:       sdk.RawData(fmt.Sprintf("%v", timestamp)),
	}
	return rec, nil
}
