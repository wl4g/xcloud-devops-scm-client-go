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
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Watcher interface {
	watch()
}

type RefresherWatcher struct {
	watchUri  string
	timeoutMs int64
	registry  *Registry
}

func (w *RefresherWatcher) Startup() {
	// Loop watching.
	for true {
		err, data := w.doWatchLongPolling()
		if err != nil {
			log.Printf("Failed to watch response. %p", err)
		} else if data != nil {
			doOnChangePropertiesSet(w.registry, data)
		}
		time.Sleep(1 * time.Second)
	}
}

func (w *RefresherWatcher) doWatchLongPolling() (error, []byte) {
	req, err := http.NewRequest("GET", w.watchUri, bytes.NewReader([]byte("")))
	if err != nil {
		return err, nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")

	// Do req.
	httpClient := &http.Client{Timeout: time.Duration(w.timeoutMs * 1000)}
	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		//log.Printf("Failed to req.")
		return err, nil
	}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Printf("Failed to get response body, %p", err)
		return err, nil
	}
	//log.Printf("Receive response message, %s", string(ret))
	return nil, ret
}

func doOnChangePropertiesSet(registry *Registry, data []byte) {
	for _, listener := range registry.Listeners() {
		listener(data)
	}
}
