// Copyright © 2022 Meroxa, Inc.
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

package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// CreatePublicationOptions contains additional options for creating a publication.
// If AllTables and Tables are both true and not empty at the same time,
// publication creation will fail.
type CreatePublicationOptions struct {
	Tables            []string
	PublicationParams []string
}

// CreatePublication creates a publication.
func CreatePublication(ctx context.Context, conn *pgconn.PgConn, name string, opts CreatePublicationOptions) error {
	if len(opts.Tables) == 0 {
		return fmt.Errorf("publication %q requires at least one table", name)
	}

	forTableString := fmt.Sprintf("FOR TABLE %s", strings.Join(opts.Tables, ", "))

	var publicationParams string
	if len(opts.PublicationParams) > 0 {
		publicationParams = fmt.Sprintf("WITH (%s)", strings.Join(opts.PublicationParams, ", "))
	}

	sql := fmt.Sprintf("CREATE PUBLICATION %q %s %s", name, forTableString, publicationParams)

	mrr := conn.Exec(ctx, sql)
	return mrr.Close()
}

// DropPublicationOptions contains additional options for dropping a publication.
type DropPublicationOptions struct {
	IfExists bool
}

// DropPublication drops a publication.
func DropPublication(ctx context.Context, conn *pgconn.PgConn, publicationName string, options DropPublicationOptions) error {
	var ifExistsString string
	if options.IfExists {
		ifExistsString = "IF EXISTS"
	}

	sql := fmt.Sprintf("DROP PUBLICATION %s %q", ifExistsString, publicationName)

	mrr := conn.Exec(ctx, sql)
	return mrr.Close()
}
