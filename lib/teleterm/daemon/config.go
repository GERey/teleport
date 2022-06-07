/*
Copyright 2021 Gravitational, Inc.

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

package daemon

import (
	"context"

	"github.com/gravitational/teleport/lib/teleterm/clusters"

	"github.com/gravitational/trace"

	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
)

type cluster interface {
	GetServers(ctx context.Context) ([]clusters.Server, error)
}

type ClusterStorer interface {
	GetByURI(uri string) (cluster, error)
}

// This is used to wrap a real clusters.Storage so it meets our interface
type StorageWrapper struct {
	*clusters.Storage
}

func (sw StorageWrapper) GetByURI(uri string) (cluster, error) {
	return sw.Storage.GetByURI(uri)
}

// Config is the cluster service config
type Config struct {
	// Storage is a storage service that reads/writes to tsh profiles
	Storage ClusterStorer
	// Clock is a clock for time-related operations
	Clock clockwork.Clock
	// InsecureSkipVerify is an option to skip HTTPS cert check
	InsecureSkipVerify bool
	// Log is a component logger
	Log *logrus.Entry
}

// CheckAndSetDefaults checks the configuration for its validity and sets default values if needed
func (c *Config) CheckAndSetDefaults() error {
	if c.Storage == nil {
		return trace.BadParameter("missing cluster storage")
	}

	if c.Clock == nil {
		c.Clock = clockwork.NewRealClock()
	}

	if c.Log == nil {
		c.Log = logrus.NewEntry(logrus.StandardLogger()).WithField(trace.Component, "daemon")
	}

	return nil
}
