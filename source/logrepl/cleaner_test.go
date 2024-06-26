// Copyright © 2024 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logrepl

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/conduitio/conduit-connector-postgres/test"
	"github.com/matryer/is"
)

func Test_Cleanup(t *testing.T) {
	conn := test.ConnectSimple(context.Background(), t, test.RepmgrConnString)

	tests := []struct {
		desc  string
		setup func(t *testing.T)
		conf  CleanupConfig

		wantErr error
	}{
		{
			desc: "drops slot and pub",
			conf: CleanupConfig{
				URL:             test.RepmgrConnString,
				SlotName:        "conduitslot1",
				PublicationName: "conduitpub1",
			},
			setup: func(t *testing.T) {
				table := test.SetupTestTable(context.Background(), t, conn)
				test.CreatePublication(t, conn, "conduitpub1", []string{table})
				test.CreateReplicationSlot(t, conn, "conduitslot1")
			},
		},
		{
			desc: "drops pub slot unspecified",
			conf: CleanupConfig{
				URL:             test.RepmgrConnString,
				PublicationName: "conduitpub2",
			},
			setup: func(t *testing.T) {
				table := test.SetupTestTable(context.Background(), t, conn)
				test.CreatePublication(t, conn, "conduitpub2", []string{table})
			},
		},
		{
			desc: "drops slot pub unspecified",
			conf: CleanupConfig{
				URL:      test.RepmgrConnString,
				SlotName: "conduitslot3",
			},
			setup: func(t *testing.T) {
				test.CreateReplicationSlot(t, conn, "conduitslot3")
			},
		},
		{
			desc: "drops pub slot missing",
			conf: CleanupConfig{
				URL:             test.RepmgrConnString,
				SlotName:        "conduitslot4",
				PublicationName: "conduitpub4",
			},
			setup: func(t *testing.T) {
				table := test.SetupTestTable(context.Background(), t, conn)
				test.CreatePublication(t, conn, "conduitpub4", []string{table})
			},
			wantErr: errors.New(`replication slot "conduitslot4" does not exist`),
		},
		{
			desc: "drops slot, pub missing", // no op
			conf: CleanupConfig{
				URL:             test.RepmgrConnString,
				SlotName:        "conduitslot5",
				PublicationName: "conduitpub5",
			},
			setup: func(t *testing.T) {
				test.CreateReplicationSlot(t, conn, "conduitslot5")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			is := is.New(t)

			if tc.setup != nil {
				tc.setup(t)
			}

			err := Cleanup(context.Background(), tc.conf)

			if tc.wantErr != nil {
				is.True(strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				is.NoErr(err)
			}
		})
	}
}
