package ttt

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestStartNewRecord(t *testing.T) {
	tdb := withDb("someRecords.csv", t)
	defer tdb.Close()
	err := tdb.StartRecord(parseTime("2022-02-02T09:17:02+01:00"))
	if err != nil {
		t.Fatal(err)
	}
	expected := withDb("someRecordsWithStartNew.csv", t)
	compareDb(expected, tdb, t)
}
func TestStartRecordAllreadyActive(t *testing.T) {
	tdb := withDb("someRecordsWithStartNew.csv", t)
	defer tdb.Close()
	err := tdb.StartRecord(parseTime("2022-02-02T09:17:02+01:00"))
	if !errors.Is(err, ActiveRecordExistsError) {
		t.Fatal("Expected ActiveRecordExistsError. Got", err)
	}
	expected := withDb("someRecordsWithStartNew.csv", t)
	compareDb(expected, tdb, t)
}

func TestEndRecord(t *testing.T) {
	tdb := withDb("someRecordsWithStartNew.csv", t)
	defer tdb.Close()
	err := tdb.EndRecord(parseTime("2022-02-02T17:59:00+01:00"))
	if err != nil {
		t.Fatal(err)
	}
	expected := withDb("someRecordsWithEnd.csv", t)
	compareDb(expected, tdb, t)
}
func TestEndRecordNotFound(t *testing.T) {
	tdb := withDb("someRecordsWithEnd.csv", t)
	defer tdb.Close()
	err := tdb.EndRecord(parseTime("2022-02-02T17:59:00+01:00"))
	if !errors.Is(err, NoActiveRecordError) {
		t.Fatal("Expected NoActiveRecordError. Got", err)
	}
	expected := withDb("someRecordsWithEnd.csv", t)
	compareDb(expected, tdb, t)
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("malformed time string " + s + err.Error())
	}
	return t
}

func withDb(csvFile string, t *testing.T) TimeTrackingDb {
	f, err := os.CreateTemp(t.TempDir(), csvFile+".ttt_test_*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	tdb, err := CreateDb(f.Name())
	if err != nil {

		fmt.Println(f.Name())
		t.Fatal(err)
	}
	cmd := exec.Command("sqlite3", "-csv", f.Name(), ".import testdata/"+csvFile+" records")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out))
	}
	return tdb
}

func compareDb(expectedDb, actualDb TimeTrackingDb, t *testing.T) {
	cmd := exec.Command("sqldiff", expectedDb.filename, actualDb.filename)
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if len(out) > 0 {
		t.Error(string(out))
	}
}
