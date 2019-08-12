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
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	regex         = regexp.MustCompile(`(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})(\.(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})){3}`)
	instanceCache *ReleaseInstance
)

/**
 * Get hardware information and process unique identification.
 */
func GetReleaseInstance(netcard string) *ReleaseInstance {
	// Get by cache.
	if instanceCache != nil {
		return instanceCache
	}

	// Default local host-name.
	hostname, err := os.Hostname()
	if err != nil {
		log.Panicf("Failed to Get local hostname. %s", err)
	}

	/*
	 * Use the specified network card name to correspond to IP.
	 */
	if common.IsEmpty(netcard) { // Compatible system environment variables.
		netcard = os.Getenv(KeyOSNetcard)
	}
	if !common.IsEmpty(netcard) {
		nis, _ := net.Interfaces()
	ok:
		for _, ni := range nis {
			if strings.EqualFold(netcard, ni.Name) {
				//fmt.Printf("Found network interfaces for - '%s'", ni.HardwareAddr)
				address, err := ni.Addrs()
				if err != nil {
					log.Panicf("Failed to Get host by addrs: %s, %s", address, err)
				}
				for _, addr := range address {
					_addr := addr.String()
					if len(regex.FindAllString(_addr, -1)) > 0 {
						a := strings.Split(_addr, "/")
						if len(a) >= 2 {
							hostname = a[0]
						} else {
							hostname = _addr
						}
						break ok
					}
				}
			}
		}
	}
	instanceCache := &ReleaseInstance{Host: hostname, Port: -1}
	return instanceCache
}

/**
 * Clear instance cache.
 */
func ClearCache() {
	instanceCache = nil
}
