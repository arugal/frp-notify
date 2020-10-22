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

package version

import (
	"encoding/json"
	"runtime"
)

var (
	Version   = ""
	GoVersion = runtime.Version()
	BuildTime = ""
	GitHash   = ""
)

type Info struct {
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
	BuildTime string `json:"buildTime"`
	GitHash   string `json:"gitHash"`
}

func Get() string {
	info := Info{
		Version:   Version,
		GoVersion: GoVersion,
		BuildTime: BuildTime,
		GitHash:   GitHash,
	}

	data, _ := json.Marshal(info)
	return string(data)
}
