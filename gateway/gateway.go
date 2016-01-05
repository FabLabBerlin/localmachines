package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var netSwitches []NetSwitch

type NetSwitch struct {
	Id        int64
	MachineId int64
	UrlOn     string
	UrlOff    string
	On        bool
}

func Init(dsn string) (err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("db connection error: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, machine_id, url_on, url_off FROM netswitch")
	if err != nil {
		return fmt.Errorf("db query: %v", err)
	}

	netSwitches = make([]NetSwitch, 0, 20)

	for rows.Next() {
		ns := NetSwitch{}
		err = rows.Scan(&ns.Id, &ns.MachineId, &ns.UrlOn, &ns.UrlOff)
		if err != nil {
			return fmt.Errorf("rows scan: %v", err)
		}
		netSwitches = append(netSwitches, ns)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("db rows err: %v", err)
	}

	return
}

func main() {
	dsn := flag.String("dsn", "root:@tcp(127.0.0.1:3306)/fabsmith", "Data Source Name: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]")
	flag.Parse()
	if err := Init(*dsn); err != nil {
		log.Fatalf("Init: %v", err)
	}
}
