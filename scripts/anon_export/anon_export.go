package main

import (
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/big"
)

var (
	db *sql.DB

	tablesToDrop = []string{
		"activations",
		"auth",
		"hexaswitch",
		"locations",
		"migrations",
		"reservations",
		"settings",
		"user",
		"user_membership",
		"user_roles",
	}

	tablesWithUid = []string{
		"coupons",
		"invoice_user_memberships",
		"invoices",
		"permission",
		"purchases",
		"user_locations",
		"user_memberships",
	}

	tablesWithLocId = []string{
		"coupons",
		"invoices",
		"machines",
		"membership",
		"monthly_earnings",
		"products",
		"purchases",
		"reservation_rules",
		"user_locations",
		"user_memberships",
	}
)

func cleanInvoiceTable() (err error) {
	_, err = db.Exec(`
UPDATE invoices
SET
fastbill_id = NULL,
fastbill_no = NULL,
canceled_fastbill_id = NULL,
canceled_fastbill_no = NULL,
customer_id = NULL,
customer_no = NULL
	`)

	return
}

func deleteNonLocation1Data() (err error) {
	for _, table := range tablesWithLocId {
		if _, err = db.Exec("DELETE FROM " + table + " WHERE location_id <> 1"); err != nil {
			return fmt.Errorf("%v: %v", table, err)
		}
	}

	return
}

var generatedIds = make(map[int64]struct{})

func randId(maxId int) (int64, error) {
	for {
		id, err := rand.Int(rand.Reader, big.NewInt(int64(maxId)))
		if err != nil {
			return 0, err
		}

		if _, duplicate := generatedIds[id.Int64()]; !duplicate {
			generatedIds[id.Int64()] = struct{}{}

			return id.Int64(), nil
		}
	}
}

func newUserIds() (uidMapping map[int64]int64, err error) {
	uidMapping = make(map[int64]int64)

	rows, err := db.Query("SELECT id FROM user")
	if err != nil {
		return
	}
	defer rows.Close()

	ids := make([]int64, 0, 2000)

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		if uidMapping[id], err = randId(2 * len(ids)); err != nil {
			return nil, fmt.Errorf("rand id: %v", err)
		}
	}

	return
}

func remapUserIds(uidMapping map[int64]int64, offset int64) (err error) {
	for _, table := range tablesWithUid {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin %v: %v", table, err)
		}

		for from := range uidMapping {
			if _, err := db.Exec("UPDATE "+table+" SET user_id = ? WHERE user_id = ?", offset+from, from); err != nil {
				tx.Rollback()
				return fmt.Errorf("update 1 %v: %v", table, err)
			}
		}

		for from, to := range uidMapping {
			if _, err := db.Exec("UPDATE "+table+" SET user_id = ? WHERE user_id = ?", to, from+offset); err != nil {
				tx.Rollback()
				return fmt.Errorf("update 2 %v: %v", table, err)
			}
		}

		if tx.Commit(); err != nil {
			return fmt.Errorf("commit %v: %v", table, err)
		}
	}

	return
}

func dropUnusedTables() (err error) {
	for _, table := range tablesToDrop {
		if _, err := db.Exec("DROP TABLE " + table); err != nil {
			return fmt.Errorf("%v: %v", table, err)
		}
	}

	return
}

func Main(dbName, user, pass string) (err error) {
	if db, err = sql.Open("mysql", user+":"+pass+"@/"+dbName); err != nil {
		return
	}

	db.Exec("ALTER TABLE user_locations DROP INDEX unique_user_locations")

	if err := deleteNonLocation1Data(); err != nil {
		return fmt.Errorf("delete non location 1 data: %v", err)
	}

	uidMapping, err := newUserIds()
	if err != nil {
		return fmt.Errorf("new user ids: %v", err)
	}

	if err := remapUserIds(uidMapping, int64(1000000+10*len(uidMapping))); err != nil {
		return fmt.Errorf("remap user ids: %v", err)
	}

	if err := dropUnusedTables(); err != nil {
		return fmt.Errorf("drop unused tables: %v", err)
	}

	if err := cleanInvoiceTable(); err != nil {
		return fmt.Errorf("clean invoice table: %v", err)
	}

	return
}

func main() {
	dbName := flag.String("dbName", "fabsmith", "DB to anonymize")
	user := flag.String("user", "user", "user")
	pass := flag.String("pass", "pass", "password")
	flag.Parse()

	if err := Main(*dbName, *user, *pass); err != nil {
		log.Fatalf("%v", err)
	}
}
