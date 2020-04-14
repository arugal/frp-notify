// Copyright 2020 arugal, zhangwei24@apache.org
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

package main

import (
	"flag"
)

type Config struct {
	Addr          string
	GotifyAddress string
	GotifyToken   string
}

func (c *Config) addFlags() {
	flag.StringVar(&c.Addr, "addr", c.Addr, "server address")
	flag.StringVar(&c.GotifyAddress, "gotify-addr", c.GotifyAddress, "gotify address")
	flag.StringVar(&c.GotifyToken, "gotify-token", c.GotifyToken, "gotify token")
}
