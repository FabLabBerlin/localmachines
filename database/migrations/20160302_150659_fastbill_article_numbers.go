package main

import (
	"fmt"
	"github.com/astaxie/beego/migration"
	"strconv"
)

const (
	PAYG_EPILOG                  = 1
	PAYG_VINYLCUTTER             = 2
	PAYG_3DPRINTING              = 3
	PAYG_CNC                     = 4
	MEMBERSHIP_3DPRINTCLUB       = 6
	MEMBERSHIP_FLB_BASIC         = 7
	MEMBERSHIP_FLB_PLUS          = 8
	PAYG_OBJET24_DIMENSION_ELITE = 9
	PAYG_TROTEC                  = 12
	RESERVATION_TROTEC           = 16
	RESERVATION_EPILOG           = 17
	RESERVATION_CNC              = 18
	PERSONAL_TUTOR               = 31
	MEMBERSHIP_RnD               = 40

	NOT_ASSIGNED = -1
)

const (
	MID_CNC    = 11
	MID_EPILOG = 3
	MID_TROTEC = 17
)

var machineIdToArticleNo = map[int64]int64{
	2:          PAYG_3DPRINTING,              //  3D Printer - 1 Vincent Vega (Replicator 2)
	MID_EPILOG: PAYG_EPILOG,                  //  Laser Cutter - Epilog Zing 6030
	4:          PAYG_3DPRINTING,              //  3D Printer - 2 Jules (Replicator 5 gen)
	6:          PAYG_3DPRINTING,              //  3D Printer - 7 Fabienne (i3 Berlin)
	7:          PAYG_3DPRINTING,              //  3D Printer - 4 Mia (Replicator 5 gen)
	8:          PAYG_3DPRINTING,              //  3D Printer - 5 Pumpkin (I3 Berlin)
	9:          NOT_ASSIGNED,                 //  Electronics Desk
	10:         PAYG_3DPRINTING,              //  3D Printer - 6 Honey Bunny (I3 Berlin)
	MID_CNC:    PAYG_CNC,                     //  CNC Router
	12:         PAYG_3DPRINTING,              //  3D Printer - 3 Mr. Wallace (Replicator Z18)
	13:         PAYG_3DPRINTING,              //  3D Printer - 8 Butch (Replicator 2)
	14:         PAYG_VINYLCUTTER,             //  Vinyl Cutter
	15:         NOT_ASSIGNED,                 //  Heat Press
	16:         NOT_ASSIGNED,                 //  Table Saw
	MID_TROTEC: PAYG_TROTEC,                  //  Laser Cutter - Trotec
	18:         PAYG_CNC,                     //  PCB CNC Machine (LPKF)
	19:         PERSONAL_TUTOR,               //  Tutor - Max
	20:         PERSONAL_TUTOR,               //  Tutor - Ahmad
	21:         PERSONAL_TUTOR,               //  Tutor - Laszlo
	23:         PERSONAL_TUTOR,               //  Tutor - Erich
	24:         PERSONAL_TUTOR,               //  Tutor - Jens
	25:         PERSONAL_TUTOR,               //  Tutor - Stefan
	26:         PERSONAL_TUTOR,               //  Tutor - Madeleine
	27:         PERSONAL_TUTOR,               //  Tutor - Yair
	28:         NOT_ASSIGNED,                 //  Form 1+
	29:         PAYG_3DPRINTING,              //  3D Printer - 1mm
	31:         PAYG_3DPRINTING,              //  3D Printer - Replicator Mini
	32:         NOT_ASSIGNED,                 //  Drill Press (Flott)
	33:         NOT_ASSIGNED,                 //  Router Table
	34:         PAYG_OBJET24_DIMENSION_ELITE, //  3D Printer - Objet24
	35:         PAYG_OBJET24_DIMENSION_ELITE, //  3D Printer - Dimension Elite
	37:         NOT_ASSIGNED,                 //  Big Rep
	38:         NOT_ASSIGNED,                 //  Textima Industrial Sewing Machine
	39:         NOT_ASSIGNED,                 //  iPad Structure Sensor & Kinect
	40:         PAYG_3DPRINTING,              //  3D Printer - 9 Angelo (Replicator 2)
	41:         NOT_ASSIGNED,                 //  Mauro Test Machine
	42:         NOT_ASSIGNED,                 //  Brother Knitting Machine
	43:         NOT_ASSIGNED,                 //  Desktop Sewing Machine
	44:         NOT_ASSIGNED,                 //  Overlock
}

var reservationNames = map[int]string{
	MID_CNC:    "Maschinenreservierung (CNC)",
	MID_EPILOG: "Maschinenreservierung (Epilog Lasercutter)",
	MID_TROTEC: "Maschinenreservierung (Trotec Lasercutter)",
}

// DO NOT MODIFY
type FastbillArticleNumbers_20160302_150659 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &FastbillArticleNumbers_20160302_150659{}
	m.Created = "20160302_150659"
	migration.Register("FastbillArticleNumbers_20160302_150659", m)
}

// Run the migrations
func (m *FastbillArticleNumbers_20160302_150659) Up() {
	m.SQL("ALTER TABLE products ADD COLUMN machine_id INT(11) AFTER user_id")
	m.SQL("ALTER TABLE products ADD COLUMN fastbill_article_no INT(11) AFTER machine_id")
	m.SQL(fmt.Sprintf("UPDATE products SET fastbill_article_no = %v WHERE type = 'tutor'",
		PERSONAL_TUTOR))
	// Machine products (-> article numbers)
	m.SQL(`
		INSERT INTO products 
		            (location_id, 
		             type, 
		             machine_id, 
		             name, 
		             price, 
		             price_unit) 
		SELECT location_id, 
		       'machine', 
		       id, 
		       name, 
		       price, 
		       price_unit 
		FROM   machines
		WHERE location_id = 1
	`)
	// Reservation products (-> article numbers)
	for machineId, name := range reservationNames {
		m.SQL(`
			INSERT INTO products 
			            (location_id, 
			             type, 
			             machine_id, 
			             name, 
			             price, 
			             price_unit) 
			SELECT location_id, 
			       'reservation', 
			       id, 
			       '` + name + `', 
			       price, 
			       price_unit 
			FROM   machines
			WHERE  id = ` + strconv.Itoa(machineId))
	}

	for machineId, articleNo := range machineIdToArticleNo {
		if articleNo != NOT_ASSIGNED {
			cmd := fmt.Sprintf("UPDATE products SET fastbill_article_no = %v WHERE machine_id = %v and type = 'machine'",
				articleNo, machineId)
			m.SQL(cmd)
		}
	}
}

// Reverse the migrations
func (m *FastbillArticleNumbers_20160302_150659) Down() {
	m.SQL("ALTER TABLE products DROP COLUMN machine_id")
	m.SQL("ALTER TABLE products DROP COLUMN fastbill_article_no")
	m.SQL("DELETE FROM products WHERE type = 'machine'")
	m.SQL("DELETE FROM products WHERE type = 'reservation'")
}
