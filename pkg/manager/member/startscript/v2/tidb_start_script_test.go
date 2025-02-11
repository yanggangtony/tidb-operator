// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"testing"

	"github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega"
)

func TestRenderTiDBStartScript(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type testcase struct {
		name string

		modifyTC     func(tc *v1alpha1.TidbCluster)
		expectScript string
	}

	cases := []testcase{
		{
			name:     "basic",
			modifyTC: func(tc *v1alpha1.TidbCluster) {},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=start-script-test-pd:2379 \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "with PDAddresses but without PreferPDAddressesOverDiscovery",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.PDAddresses = []string{"${PD_DOMAIN}:2380", "another.pd:2380"}
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=start-script-test-pd:2379 \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "with PDAddresses and PreferPDAddressesOverDiscovery",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.PDAddresses = []string{"${PD_DOMAIN}:2380", "another.pd:2380"}
				tc.Spec.StartScriptV2FeatureFlags = []v1alpha1.StartScriptV2FeatureFlag{
					v1alpha1.StartScriptV2FeatureFlagPreferPDAddressesOverDiscovery,
				}
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=${PD_DOMAIN}:2380,another.pd:2380 \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "set plugin",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.TiDB.Plugins = []string{"plugin-1", "plugin-2"}
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=start-script-test-pd:2379 \
--config=/etc/tidb/tidb.toml"
ARGS="${ARGS} --plugin-dir=/plugins --plugin-load=plugin-1,plugin-2"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "enable tidb binlog",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				enable := true
				tc.Spec.TiDB.BinlogEnabled = &enable
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=start-script-test-pd:2379 \
--config=/etc/tidb/tidb.toml"
ARGS="${ARGS} --enable-binlog=true"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "enable tidb binlog when pump is setted",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.TiDB.BinlogEnabled = nil
				tc.Spec.Pump = &v1alpha1.PumpSpec{}
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=start-script-test-pd:2379 \
--config=/etc/tidb/tidb.toml"
ARGS="${ARGS} --enable-binlog=true"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "non-empty cluster domain",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.ClusterDomain = "test.com"
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc.test.com \
--host=0.0.0.0 \
--path=start-script-test-pd:2379 \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "across k8s with setting cluster domain",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.ClusterDomain = "test.com"
				tc.Spec.AcrossK8s = true
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}
pd_url=start-script-test-pd:2379
encoded_domain_url=$(echo $pd_url | base64 | tr "\n" " " | sed "s/ //g")
discovery_url=start-script-test-discovery.start-script-test-ns:10261
until result=$(wget -qO- -T 3 http://${discovery_url}/verify/${encoded_domain_url} 2>/dev/null | sed 's/http:\/\///g'); do
    echo "waiting for the verification of PD endpoints ..."
    sleep $((RANDOM % 5))
done

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc.test.com \
--host=0.0.0.0 \
--path=${result} \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "across k8s without setting cluster domain",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.ClusterDomain = ""
				tc.Spec.AcrossK8s = true
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}
pd_url=start-script-test-pd:2379
encoded_domain_url=$(echo $pd_url | base64 | tr "\n" " " | sed "s/ //g")
discovery_url=start-script-test-discovery.start-script-test-ns:10261
until result=$(wget -qO- -T 3 http://${discovery_url}/verify/${encoded_domain_url} 2>/dev/null | sed 's/http:\/\///g'); do
    echo "waiting for the verification of PD endpoints ..."
    sleep $((RANDOM % 5))
done

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=${result} \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "heterogeneous cluster without local pd",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.PD = nil
				tc.Spec.Cluster = &v1alpha1.TidbClusterRef{Name: "target-cluster"}
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=target-cluster-pd:2379 \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
		{
			name: "heterogeneous cluster when across k8s",
			modifyTC: func(tc *v1alpha1.TidbCluster) {
				tc.Spec.PD = nil
				tc.Spec.Cluster = &v1alpha1.TidbClusterRef{Name: "target-cluster"}
				tc.Spec.AcrossK8s = true
			},
			expectScript: `#!/bin/sh

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"
if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

TIDB_POD_NAME=${POD_NAME:-$HOSTNAME}
pd_url=start-script-test-pd:2379
encoded_domain_url=$(echo $pd_url | base64 | tr "\n" " " | sed "s/ //g")
discovery_url=start-script-test-discovery.start-script-test-ns:10261
until result=$(wget -qO- -T 3 http://${discovery_url}/verify/${encoded_domain_url} 2>/dev/null | sed 's/http:\/\///g'); do
    echo "waiting for the verification of PD endpoints ..."
    sleep $((RANDOM % 5))
done

ARGS="--store=tikv \
--advertise-address=${TIDB_POD_NAME}.start-script-test-tidb-peer.start-script-test-ns.svc \
--host=0.0.0.0 \
--path=${result} \
--config=/etc/tidb/tidb.toml"

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`,
		},
	}

	for _, c := range cases {
		t.Logf("test case: %s", c.name)

		tc := &v1alpha1.TidbCluster{
			Spec: v1alpha1.TidbClusterSpec{
				TiDB: &v1alpha1.TiDBSpec{},
			},
		}
		tc.Name = "start-script-test"
		tc.Namespace = "start-script-test-ns"

		if c.modifyTC != nil {
			c.modifyTC(tc)
		}

		script, err := RenderTiDBStartScript(tc)
		g.Expect(err).Should(gomega.Succeed())
		if diff := cmp.Diff(c.expectScript, script); diff != "" {
			t.Errorf("unexpected (-want, +got): %s", diff)
		}
		g.Expect(validateScript(script)).Should(gomega.Succeed())
	}
}
