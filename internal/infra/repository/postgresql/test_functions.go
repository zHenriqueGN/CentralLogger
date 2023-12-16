package postgresql

import (
	"database/sql"
	"testing"
)

func resetTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS logs;")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS systems;")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE systems (
    		id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY,
        	name VARCHAR(50) NOT NULL,
        	description VARCHAR(200) NOT NULL,
        	version VARCHAR(15) NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE logs (
        	id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY,
        	system_id VARCHAR(36) NOT NULL REFERENCES systems (id),
        	level VARCHAR(30) NOT NULL,
        	status VARCHAR(30) NOT NULL,
        	message VARCHAR NOT NULL,
        	time_stamp TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
}
