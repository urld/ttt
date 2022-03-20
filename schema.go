// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"database/sql"
)

type schemaPatch []string

func initSchema(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range baseSchema() {
		_, err := tx.Exec(stmt)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func applySchema(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	revRow := tx.QueryRow(selectRevisionStmt)
	var dbRev int
	err = revRow.Scan(&dbRev)
	if err != nil {
		return err
	}

	for rev, patch := range schemaPatches() {
		if rev > dbRev {
			for _, stmt := range patch {
				_, err := tx.Exec(stmt)
				if err != nil {
					return err
				}
			}
			_, err = tx.Exec(updateRevisionStmt, rev)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

const (
	selectRevisionStmt = "SELECT revision FROM schema_version;"
	updateRevisionStmt = "UPDATE schema_version SET revision=?, dts=datetime('now');"
)

func baseSchema() schemaPatch {
	return schemaPatch{
		"CREATE TABLE schema_version (revision int, dts datetime);",
		"INSERT INTO schema_version VALUES(-1, datetime('now'));",
	}

}

func schemaPatches() []schemaPatch {
	return []schemaPatch{
		{
			"CREATE TABLE config (property text, value text);",
			"INSERT INTO config VALUES('break_threshold', '6h'), ('break_deduction', '30m');",
			"INSERT INTO config VALUES('monday_hours', '7h 42m');",
			"INSERT INTO config VALUES('tuesday_hours', '7h 42m');",
			"INSERT INTO config VALUES('wednesday_hours', '7h 42m');",
			"INSERT INTO config VALUES('thursday_hours', '7h 42m');",
			"INSERT INTO config VALUES('friday_hours', '7h 42m');",
			"INSERT INTO config VALUES('saturday_hours', '0h');",
			"INSERT INTO config VALUES('sunday_hours', '0h');",
			"INSERT INTO config VALUES('holidays', 'AT');",

			"CREATE TABLE records (start datetime, end datetime);",
		},
	}

}
