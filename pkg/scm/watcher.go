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
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RefreshWatcher struct {
	serverUri string
	timeoutMs int64
	registry  *Registry
}

func (_self *RefreshWatcher) Sync() {
	select {}
}

func (_self *RefreshWatcher) Startup() *RefreshWatcher {
	// Loop watching.
	go func() {
		for true {
			err, meta := _self.createWatchLongPolling()
			if err != nil {
				log.Printf("Failed to watching. %s", err.Error())
			} else if meta != nil {
				_self.doOnChangedProperties(_self.registry, meta)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return _self
}

func (_self *RefreshWatcher) createWatchLongPolling() (error, *ReleaseMeta) {
	watchUrl := _self.serverUri + "/watch"
	err, resp, data := _self.doExchange(watchUrl, "", "GET", _self.timeoutMs)
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

func (_self *RefreshWatcher) doExchange(url string, params string, method string, timeoutMs int64) (error, *http.Response, []byte) {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(params)))
	if err != nil {
		return err, nil, nil
	}
	_self.addHeader(req)

	// Do req.
	httpClient := &http.Client{Timeout: time.Duration(timeoutMs * 1000)}
	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		//log.Printf("Failed to req.")
		return err, resp, nil
	}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Printf("Failed to get response body, %p", err)
		return err, resp, nil
	}
	//log.Printf("Receive response message, %s", string(ret))
	return nil, resp, ret
}

func (_self *RefreshWatcher) addHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")
}

func (_self *RefreshWatcher) doOnChangedProperties(registry *Registry, meta *ReleaseMeta) {
	// Fetching release sources.
	fetchUrl := _self.serverUri + "/source"
	err, resp, data := _self.doExchange(fetchUrl, "", "GET", 4000)
	if err != nil {
		log.Printf("Failed to fetch property soruces. %s", err)
		return
	}

	// Extract property sources.
	msgResp := ReleaseMessageResp{}
	err = jsoniter.Unmarshal(data, resp)
	if err != nil {
		log.Printf("Failed to extract property soruces. %s", err)
	}
	releaseMessage := msgResp.data["release-source"]

	// Callback listeners.
	for _, listener := range registry.Listeners() {
		listener(meta, releaseMessage)
	}
}
