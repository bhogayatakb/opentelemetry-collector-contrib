// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vcenterreceiver // import github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestGetClusters(t *testing.T) {
	simulator.Test(func(ctx context.Context, c *vim25.Client) {
		finder := find.NewFinder(c)
		client := vcenterClient{
			vimDriver: c,
			finder:    finder,
		}
		dc, err := finder.DefaultDatacenter(ctx)
		require.NoError(t, err)
		clusters, err := client.Clusters(ctx, dc)
		require.NoError(t, err)
		require.NotEmpty(t, clusters, 0)
	})
}

func TestGetResourcePools(t *testing.T) {
	simulator.Test(func(ctx context.Context, c *vim25.Client) {
		finder := find.NewFinder(c)
		client := vcenterClient{
			vimDriver: c,
			finder:    finder,
		}
		resourcePools, err := client.ResourcePools(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, resourcePools)
	})
}

func TestGetVMs(t *testing.T) {
	simulator.Test(func(ctx context.Context, c *vim25.Client) {
		finder := find.NewFinder(c)
		client := vcenterClient{
			vimDriver: c,
			finder:    finder,
		}
		vms, err := client.VMs(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, vms)
	})
}

func TestSessionReestablish(t *testing.T) {
	simulator.Test(func(ctx context.Context, c *vim25.Client) {
		sm := session.NewManager(c)
		moClient := &govmomi.Client{
			Client:         c,
			SessionManager: sm,
		}
		pw, _ := simulator.DefaultLogin.Password()
		client := vcenterClient{
			vimDriver: c,
			cfg: &Config{
				Username: simulator.DefaultLogin.Username(),
				Password: pw,
				Endpoint: fmt.Sprintf("%s://%s", c.URL().Scheme, c.URL().Host),
				TLSClientSetting: configtls.TLSClientSetting{
					Insecure: true,
				},
			},
			moClient: moClient,
		}
		err := sm.Logout(ctx)
		require.NoError(t, err)

		connected, err := client.moClient.SessionManager.SessionIsActive(ctx)
		require.NoError(t, err)
		require.False(t, connected)

		err = client.EnsureConnection(ctx)
		require.NoError(t, err)

		connected, err = client.moClient.SessionManager.SessionIsActive(ctx)
		require.NoError(t, err)
		require.True(t, connected)
	})
}
