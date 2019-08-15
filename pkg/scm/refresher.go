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
	"github.com/wl4g/super-devops-scm-agent/pkg/common"
	"github.com/wl4g/super-devops-scm-agent/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type ConfigListener func(meta *ReleaseMeta, release ReleaseMessage)

type RefreshOption struct {
	Server     string
	Netcard    string
	Cluster    string
	Endpoint   string
	Namespaces string
}

type DefaultRefresher struct {
	RefreshOption
	registry *Registry
	watcher  *DefaultWatcher
}

func NewRefresher(opt RefreshOption) (*DefaultRefresher, error) {
	// Check options.
	checkOptions(&opt)

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
	httpClient := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}
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
	releaseI := GetReleaseInstance(_self.RefreshOption)
	get := GetRelease{Instance: *releaseI}
	get.Meta = *meta
	get.Cluster = _self.Cluster
	get.Namespaces = strings.Split(_self.Namespaces, ",")

	// Fetching release sources.
	fetchUrl := _self.Server + UriEndpointRefreshFetch
	err, resp, data := _self.doExchange(fetchUrl, get.AsJsonString(), "GET", DefaultFetchTimeout)
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
	releaseMsg := msgResp.data[KeyReleaseSource]

	// Callback listeners.
	for _, listener := range registry.Listeners() {
		listener(meta, releaseMsg)
	}
}

/**
 * Check and use default refresh options.
 */
func (_self *RefreshOption) checkAndUseDefault() {
	if common.IsEmpty(_self.Server) {
		errors.FatalExit("Illegal scm server address. %s", _self.Server)
	}
	if common.IsEmpty(_self.Netcard) {
		log.Printf("Default local hostname is used.")
	}
	if common.IsEmpty(_self.Cluster) {
		errors.FatalExit("SCM cluster must not be empty!")
	}
	if common.IsEmpty(_self.Endpoint) {
		errors.FatalExit("SCM endpoint must not be empty!")
	}
	if common.IsEmpty(_self.Namespaces) {
		errors.FatalExit("SCM namespaces must not be empty!")
	}
}

/**
 * Check refresh options.
 */
func checkOptions(opt *RefreshOption) {
	if opt == nil {
		log.Panicf("Illegal refresh option! %s", common.ToJSONString(opt))
	}
	opt.checkAndUseDefault()
}
