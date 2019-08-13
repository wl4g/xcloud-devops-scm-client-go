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
package common

import (
	"os"
	"strconv"
)

/**
 * Get string by environment
 */
func GetEnvString(key string, defaultValue string) (string, error) {
	val := os.Getenv(key)
	if IsEmpty(val) {
		val = defaultValue
	}
	return val, nil
}

/**
 * Get int by environment
 */
func GetEnvInt(key string, defaultValue int) (int, error) {
	val := os.Getenv(key)
	if !IsEmpty(val) {
		return strconv.Atoi(val)
	}
	return defaultValue, nil
}

/**
 * Get int32 by environment
 */
func GetEnvInt32(key string, defaultValue int32) (int32, error) {
	val := os.Getenv(key)
	if !IsEmpty(val) {
		value, err := strconv.ParseInt(val, 10, 32)
		return int32(value), err
	}
	return defaultValue, nil
}

/**
 * Get int64 by environment
 */
func GetEnvInt64(key string, defaultValue int64) (int64, error) {
	val := os.Getenv(key)
	if !IsEmpty(val) {
		return strconv.ParseInt(val, 10, 64)
	}
	return defaultValue, nil
}

/**
 * Get float32 by environment
 */
func GetEnvFloat32(key string, defaultValue float32) (float32, error) {
	val := os.Getenv(key)
	if !IsEmpty(val) {
		value, err := strconv.ParseFloat(val, 32)
		return float32(value), err
	}
	return defaultValue, nil
}

/**
 * Get float64 by environment
 */
func GetEnvFloat64(key string, defaultValue float64) (float64, error) {
	val := os.Getenv(key)
	if !IsEmpty(val) {
		return strconv.ParseFloat(val, 64)
	}
	return defaultValue, nil
}
