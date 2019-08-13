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

import "github.com/wl4g/super-devops-scm-agent/pkg/common"

//
// --- SCM message ---
//

type ReleaseInstance struct {
	Host     string `json:"host"`
	Endpoint string `json:"endpoint"`
}

type ReleaseMeta struct {
	Version   string `json:"version"`
	ReleaseId string `json:"releaseId"`
}

type GenericInfo struct {
	Group      string      `json:"group"`
	Namespaces []string    `json:"namespaces"`
	Meta       ReleaseMeta `json:"meta"`
}

type GetRelease struct {
	Instance ReleaseInstance `json:"instance"`
	GenericInfo
}

type ReleaseMessage struct {
	GetRelease
	PropertySources []ReleasePropertySource `json:"propertySources"`
}

type ReleasePropertySource struct {
	name   string                 `json:"name"`
	source map[string]interface{} `json:"source"`
}

//
// --- Message for function ---
//

func (msg *GenericInfo) AsJsonString() string {
	return common.ToJSONString(msg)
}

func (msg *ReleaseMeta) AsJsonString() string {
	return common.ToJSONString(msg)
}

//
// --- Response ---
//

type ReleaseMessageResp struct {
	RespBase
	data map[string]ReleaseMessage `json:"data"`
}
