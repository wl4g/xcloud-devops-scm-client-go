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
	"github.com/wl4g/super-devops-scm-agent/pkg/common"
	"sync"
)

type AlreadyRegisteredError struct {
	name     string
	listener ConfigListener
}

func (err AlreadyRegisteredError) Error() string {
	return "Duplicate config listener registration attempted"
}

type Registry struct {
	mtx       sync.RWMutex
	listeners map[string]ConfigListener
}

func newRegistry() *Registry {
	return &Registry{
		listeners: map[string]ConfigListener{},
	}
}

func (r *Registry) Register(name string, l ConfigListener) error {
	r.mtx.Lock()
	// Check exists.
	if _, exists := r.listeners[name]; exists {
		return AlreadyRegisteredError{name: name, listener: l}
	}
	r.listeners[name] = l
	r.mtx.Unlock()
	return nil
}

func (r *Registry) Unregister(name string) bool {
	r.mtx.Lock()
	// Check exists.
	if _, exists := r.listeners[name]; !exists {
		r.mtx.Unlock()
		return false
	}
	delete(r.listeners, name)
	// Confirm check.
	if _, exists := r.listeners[name]; !exists {
		r.mtx.Unlock()
		return true
	}
	r.mtx.Unlock()
	return false
}

func (r *Registry) Listeners() []ConfigListener {
	return r.GetListeners(nil)
}

func (r *Registry) GetListeners(names []string) []ConfigListener {
	var listeners []ConfigListener
	for _name, l := range r.listeners {
		if names == nil || len(names) <= 0 || common.StringsContains(names, _name) {
			listeners = append(listeners, l)
		}
	}
	return listeners
}
