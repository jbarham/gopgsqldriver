// Copyright 2011 John E. Barham. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pgsqldriver

import (
	"database/sql"
	"math"
	"testing"
)

type rec struct {
	tf  bool
	i32 int
	i64 int64
	s   string
	b   []byte
}

var testTuples = []rec{
	{false, math.MinInt32, math.MinInt64, "hello world", []byte{0xDE, 0xAD}},
	{true, math.MaxInt32, math.MaxInt64, "Γεια σας κόσμο", []byte{0xBE, 0xEF}},
}

func chkerr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestPq(t *testing.T) {
	// pgsqldriver registers itself with sql.Register using name "postgres".
	db, err := sql.Open("postgres", "dbname=testdb")
	chkerr(t, err)

	// Create test table, and schedule its deletion.
	_, err = db.Exec("CREATE TABLE gopq_test (tf bool, i32 int, i64 bigint, s text)")
	chkerr(t, err)
	defer db.Exec("DROP TABLE gopq_test")

	// Insert test rows.
	stmt, err := db.Prepare("INSERT INTO gopq_test VALUES ($1, $2, $3, $4)")
	chkerr(t, err)
	defer stmt.Close()
	for _, row := range testTuples {
		_, err = stmt.Exec(row.tf, row.i32, row.i64, row.s)
		chkerr(t, err)
	}

	// Verify that all test rows were inserted.
	rows, err := db.Query("SELECT COUNT(*) FROM gopq_test")
	chkerr(t, err)
	if !rows.Next() {
		t.Fatal("Result.Next failed")
	}
	var count int
	err = rows.Scan(&count)
	chkerr(t, err)
	if count != len(testTuples) {
		t.Fatalf("invalid row count %d, expected %d", count, len(testTuples))
	}
	rows.Close()

	// Retrieve inserted rows and verify inserted values.
	rows, err = db.Query("SELECT * FROM gopq_test")
	chkerr(t, err)
	for i := 0; rows.Next(); i++ {
		var tf bool
		var i32 int
		var i64 int64
		var s string

		err := rows.Scan(&tf, &i32, &i64, &s)
		if err != nil {
			t.Fatal("scan error:", err)
		}
		if tf != testTuples[i].tf {
			t.Fatal("bad bool")
		}
		if i32 != testTuples[i].i32 {
			t.Fatal("bad int32")
		}
		if i64 != testTuples[i].i64 {
			t.Fatal("bad int64")
		}
		if s != testTuples[i].s {
			t.Fatal("bad string")
		}
	}
	rows.Close()
}
