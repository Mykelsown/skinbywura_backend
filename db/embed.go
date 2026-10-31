// Package db embeds the SQL migrations so the single compiled binary can migrate its own database.
package db

import "embed"

// Migrations holds every *.up.sql / *.down.sql file in db/migrations.
//
//go:embed migrations/*.sql
var Migrations embed.FS
