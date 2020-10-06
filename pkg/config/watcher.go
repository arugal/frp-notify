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

package config

import (
	"time"

	"istio.io/pkg/filewatcher"
)

var (
	watchers []WatchNotifyConfigFunc
)

type WatchNotifyConfigFunc func(cfg FRPNotifyConfig)

func RegisterConfigListener(watcher WatchNotifyConfigFunc) {
	watchers = append(watchers, watcher)
}

type FRPNotifyConfigController struct {
	bindAddress    string
	windowInterval int64

	configPath string
}

func NewConfigController(bindAddress string, windowInterval int64, configPath string) *FRPNotifyConfigController {
	return &FRPNotifyConfigController{
		bindAddress:    bindAddress,
		windowInterval: windowInterval,
		configPath:     configPath,
	}
}

func (c *FRPNotifyConfigController) Start(stop chan struct{}) {
	c.reload()

	// watcher config file
	fileWatcher := filewatcher.NewWatcher()
	defer fileWatcher.Close()
	addFileWatcher(fileWatcher, c.configPath, func() {
		c.reload()
	})
	<-stop
}

func (c *FRPNotifyConfigController) reload() {
	cfg := Load(c.configPath)
	cfg.BindAddress = c.bindAddress
	cfg.WindowInterval = c.windowInterval
	log.Debugf("config: %v", cfg)
	c.configEvent(*cfg)
}

func (c *FRPNotifyConfigController) configEvent(cfg FRPNotifyConfig) {
	for _, watcher := range watchers {
		watcher(cfg)
	}
}

// Add to the FileWatcher the provided file and execute the provided function
// on any change event for this file.
// Using a debouncing mechanism to avoid calling the callback multiple times
// per event.
func addFileWatcher(fileWatcher filewatcher.FileWatcher, file string, callback func()) {
	_ = fileWatcher.Add(file)
	go func() {
		var timerC <-chan time.Time
		for {
			select {
			case <-timerC:
				timerC = nil
				callback()
			case <-fileWatcher.Events(file):
				// Use a timer to debounce configuration updates
				if timerC == nil {
					timerC = time.After(100 * time.Millisecond)
				}
			}
		}
	}()
}
