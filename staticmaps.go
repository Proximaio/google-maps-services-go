// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// More information about Google Distance Matrix API is available on
// https://developers.google.com/maps/documentation/distancematrix/

package maps

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"golang.org/x/net/context"
)

var staticmapsAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/staticmap",
	acceptsClientID: true,
}

// StaticMap makes a Static Map API request
func (c *Client) StaticMap(ctx context.Context, r *StaticMapsRequest) (io.ReadCloser, error) {
	if len(r.Size) != 2 {
		return nil, errors.New("maps: size parameter set incorrectly")
	}

	response, err := c.getBinary(ctx, staticmapsAPI, r)
	if err != nil {
		return nil, err
	}
	if response.statusCode != 200 {
		str, err := ioutil.ReadAll(response.data)

		if err != nil {
			return nil, err
		}
		defer response.data.Close()
		return nil, errors.New(string(str))
	}

	return response.data, nil
}

func (r *StaticMapsRequest) params() url.Values {
	q := make(url.Values)

	if len(r.Size) == 2 {
		q.Set("size", fmt.Sprintf("%dx%d", r.Size[0], r.Size[1]))
	}

	if len(r.Markers.Locations) > 0 {
		q.Set("markers", fmt.Sprintf("%v|%v", r.Markers.Style, r.Markers.Locations[0]))
	}

	return q
}

// StaticMapsRequest is the request structure for Static Maps API
type StaticMapsRequest struct {
	// Size defines the rectangular dimensions of the map image, [WIDTH, HEIGHT]
	Size []uint

	// Markers defines a set of one or more markers at a set of locations
	Markers Markers
}

type Markers struct {
	Style     string
	Locations []string
}
