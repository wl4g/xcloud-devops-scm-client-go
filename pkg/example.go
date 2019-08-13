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
package main

import (
	"github.com/wl4g/super-devops-scm-agent/pkg/scm"
	"log"
)

func main() {
	log.Printf("SCM example starting...")

	opt := scm.RefreshOption{
		Server:       "http://localhost:14043/scm-server",
		Cluster:      "ExampleCollector",
		Port:         1025,
		InstanceAddr: "ExampleCollector1",
		Namespaces:   "scm-agent.yml",
	}

	// Create refresher.
	refresher, err := scm.NewRefresher(opt)
	if err != nil {
		log.Panicf("Failed to create refresher. %s", err)
	}

	// Register listener.
	_ = refresher.Registry().Register("exampleListener", func(meta *scm.ReleaseMeta, release scm.ReleaseMessage) {
		log.Printf("On change config ... for meta: %s => %s", meta.AsJsonString(), release.AsJsonString())
	})

	// Startup & watching.
	refresher.Watcher().Startup().Sync()
}
