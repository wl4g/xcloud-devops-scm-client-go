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

type ConfigListener func(meta *ReleaseMeta, release ReleaseMessage)

type Refresher struct {
	registry *Registry
	watcher  *RefreshWatcher
}

func (r *Refresher) Registry() *Registry {
	return r.registry
}

func (r *Refresher) Watcher() *RefreshWatcher {
	return r.watcher
}

func NewRefresher(watchUri string, timeoutMs int64) (*Refresher, error) {
	// New registry.
	registry := newRegistry()

	// Create watch.
	watcher := &RefreshWatcher{serverUri: watchUri, timeoutMs: timeoutMs, registry: registry}

	// Refresher.
	return &Refresher{
		registry: registry,
		watcher:  watcher,
	}, nil
}
