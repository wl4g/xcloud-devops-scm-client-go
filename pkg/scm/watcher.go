/**
 * Copyright 2017 ~ 2025 the original author or authors[983708408@qq.com].
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package scm

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
)

type DefaultWatcher struct {
	refresher DefaultRefresher
	timeoutMs int64
	registry  *Registry
}

func (_self *DefaultWatcher) Sync() {
	select {}
}

func (_self *DefaultWatcher) Startup() *DefaultWatcher {
	// Loop watching.
	go func() {
		for true {
			err, meta := _self.createWatchLongPolling()
			if err != nil {
				log.Printf("Failed to watching. %s", err.Error())
			} else if meta != nil {
				_self.refresher.refresh(_self.registry, meta)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return _self
}

func (_self *DefaultWatcher) createWatchLongPolling() (error, *ReleaseMeta) {
	watchUrl := _self.refresher.serverUri + "/watch"
	err, resp, data := _self.refresher.doExchange(watchUrl, "", "GET", _self.timeoutMs)
	if err != nil {
		return err, nil
	}

	// Update watching state
	switch resp.StatusCode {
	case 200: // On change
		meta := ReleaseMeta{}
		err = jsoniter.Unmarshal(data, meta)
		if err != nil {
			return err, nil
		}
		return nil, &meta
	case 103:
		// Not implemented
	case 304:
		break // Not modified(next long-polling)
	default:
		return &IllegalWatchExceptionError{StatusCode: resp.StatusCode}, nil
	}
	return nil, nil
}
