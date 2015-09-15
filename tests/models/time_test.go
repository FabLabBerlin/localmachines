package modelTest

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// Testing whether the time is stored and read back
// if using Beego ORM models
func TestTime(t *testing.T) {
	Convey("Testing Time", t, func() {

		myTime := TimeTest{}

		// Create a table for testing
		sql := fmt.Sprintf("CREATE TABLE %s "+
			"(id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, "+
			"time DATETIME)", myTime.TableName())
		o := orm.NewOrm()
		_, err := o.Raw(sql).Exec()

		Convey("Creating test table in the db", func() {
			So(err, ShouldBeNil)
		})

		// Get current time
		currentTime := time.Now()

		// Insert time as UTC time
		myTime.Time = currentTime
		var id int64
		id, err = o.Insert(&myTime)

		// Read time from db
		retTime := TimeTest{}
		retTime.Id = id
		errOnRet := o.Read(&retTime)

		Convey("Saving current time in to the database", func() {
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
		})

		Convey("Retrieving time from the database", func() {
			So(errOnRet, ShouldBeNil)
		})

		Convey("Compare the time that was inserted with the time read", func() {
			So(retTime.Time, ShouldHappenWithin,
				time.Duration(1)*time.Second, currentTime)
		})

		// Remove table from the database
		sql = fmt.Sprintf("DROP TABLE %s", myTime.TableName())
		_, err = o.Raw(sql).Exec()
		Convey("Removing table from the database", func() {
			So(err, ShouldBeNil)
		})
	})
}
