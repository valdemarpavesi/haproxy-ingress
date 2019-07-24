/*
Copyright 2019 The HAProxy Ingress Controller Authors.

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

package types

import (
	"fmt"
	"sort"
	"strings"
)

// NewBackendPaths ...
func NewBackendPaths(paths ...*BackendPath) BackendPaths {
	b := BackendPaths{}
	b.Add(paths...)
	return b
}

// FindEndpoint ...
func (b *Backend) FindEndpoint(target string) *Endpoint {
	for _, endpoint := range b.Endpoints {
		if endpoint.Target == target {
			return endpoint
		}
	}
	return nil
}

// AcquireEndpoint ...
func (b *Backend) AcquireEndpoint(ip string, port int, targetRef string) *Endpoint {
	endpoint := b.FindEndpoint(fmt.Sprintf("%s:%d", ip, port))
	if endpoint != nil {
		return endpoint
	}
	return b.addEndpoint(ip, port, targetRef)
}

// AddEmptyEndpoint ...
func (b *Backend) AddEmptyEndpoint() *Endpoint {
	endpoint := b.addEndpoint("127.0.0.1", 1023, "")
	endpoint.Enabled = false
	endpoint.Weight = 0
	return endpoint
}

func (b *Backend) addEndpoint(ip string, port int, targetRef string) *Endpoint {
	endpoint := &Endpoint{
		Name:      fmt.Sprintf("srv%03d", len(b.Endpoints)+1),
		IP:        ip,
		Port:      port,
		Target:    fmt.Sprintf("%s:%d", ip, port),
		Enabled:   true,
		TargetRef: targetRef,
		Weight:    1,
	}
	b.Endpoints = append(b.Endpoints, endpoint)
	return endpoint
}

// FindHostPath ...
func (b *Backend) FindHostPath(hostpath string) *BackendPath {
	for _, p := range b.Paths {
		if p.Hostpath == hostpath {
			return p
		}
	}
	return nil
}

// AddHostPath ...
func (b *Backend) AddHostPath(hostname, path string) *BackendPath {
	hostpath := hostname + path
	// add only unique paths
	backendPath := b.FindHostPath(hostpath)
	if backendPath != nil {
		return backendPath
	}
	// host's paths that references this backend
	// used on RewriteURL config
	backendPath = &BackendPath{
		ID:       fmt.Sprintf("path%02d", len(b.Paths)+1),
		Hostpath: hostpath,
		Path:     path,
	}
	b.Paths = append(b.Paths, backendPath)
	// reverse order in order to avoid overlap of sub-paths
	sort.Slice(b.Paths, func(i, j int) bool {
		return b.Paths[i].Hostpath > b.Paths[j].Hostpath
	})
	return backendPath
}

// NeedACL ...
func (b *Backend) NeedACL() bool {
	return len(b.HSTS) > 1 ||
		len(b.ProxyBodySize) > 1 || len(b.RewriteURL) > 1 || len(b.WhitelistHTTP) > 1
}

// Has ...
func (p *BackendPaths) Has(path string) bool {
	for _, item := range p.Items {
		if item.Path == path {
			return true
		}
	}
	return false
}

// IDList ...
func (p *BackendPaths) IDList() string {
	ids := make([]string, len(p.Items))
	for i, item := range p.Items {
		ids[i] = item.ID
	}
	return strings.Join(ids, " ")
}

// Add ...
func (p *BackendPaths) Add(paths ...*BackendPath) {
	for _, path := range paths {
		if path == nil {
			panic("path cannot be nil")
		}
		p.Items = append(p.Items, path)
	}
	sort.SliceStable(p.Items, func(i, j int) bool {
		return p.Items[i].Hostpath < p.Items[j].Hostpath
	})
}

// String ...
func (p *BackendPath) String() string {
	return fmt.Sprintf("%+v", *p)
}

// String ...
func (b *BackendConfigStr) String() string {
	return fmt.Sprintf("%+v", *b)
}

// String ...
func (b *BackendConfigHSTS) String() string {
	return fmt.Sprintf("%+v", *b)
}

// String ...
func (b *BackendConfigWhitelist) String() string {
	return fmt.Sprintf("%+v", *b)
}
