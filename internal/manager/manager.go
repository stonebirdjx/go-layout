/*
Copyright 2026 stonebirdjx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package manager provides a framework for managing the lifecycle of the application.
// It is inspired by controller-runtime and provides a simple and extensible way to manage the application lifecycle.
package manager

import (
	"context"
)

type Manager interface {
	Validate() error
	Start(ctx context.Context) error
}

func NewManager(ops ...Option) Manager {
	c := &controller{}
	for _, op := range ops {
		op(c)
	}
	return c
}
