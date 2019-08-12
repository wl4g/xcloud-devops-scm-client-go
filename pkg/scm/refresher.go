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

type ConfigListener func(meta *ReleaseMeta, release ReleaseMessage)

type RefreshOption struct {
	ServerUri  string
	TimeoutMs  int64
	Netcard    string
	Group      string
	Port       int
	Namespaces []string
}

type DefaultRefresher struct {
	RefreshOption
	registry *Registry
	watcher  *DefaultWatcher
}

func NewRefresher(opt RefreshOption) (*DefaultRefresher, error) {
	// New registry.
	registry := newRegistry()

	// Create refresher.
	refresher := DefaultRefresher{registry: registry}
	refresher.RefreshOption = opt

	// Create watcher.
	watcher := &DefaultWatcher{refresher: refresher, registry: registry}
	refresher.watcher = watcher

	// RefreshOption.
	return &refresher, nil
}

func (_self *DefaultRefresher) Registry() *Registry {
	return _self.registry
}

func (_self *DefaultRefresher) Watcher() *DefaultWatcher {
	return _self.watcher
}

func (_self *DefaultRefresher) doExchange(url string, params string, method string, timeoutMs int64) (error, *http.Response, []byte) {
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

func (_self *DefaultRefresher) addHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")
}

func (_self *DefaultRefresher) refresh(registry *Registry, meta *ReleaseMeta) {
	// Get release parameters.
	releaseI := ReleaseInstance{Host: GetHostAddr(_self.Netcard), Port: -1}

	get := GetRelease{Instance: releaseI}
	get.Meta = *meta
	get.Group = _self.Group
	get.Namespaces = _self.Namespaces

	// Fetching release sources.
	fetchUrl := _self.ServerUri + UriEndpointRefreshFetch
	err, resp, data := _self.doExchange(fetchUrl, get.AsJsonString(), "GET", 4000)
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
