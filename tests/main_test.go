package tests

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOnePlusOne(t *testing.T) {
	Convey("Testing 1+1", t, func() {
		res := 1 + 1
		Convey("Should equal 2", func() {
			So(res, ShouldEqual, 2)
		})
	})
}

// Fibonacci function
func fibo(rank uint64) uint64 {
	if rank == 0 {
		return 0
	} else if rank <= 2 {
		return 1
	} else {
		return fibo(rank-1) + fibo(rank-2)
	}
}

// Test example for fibonacci
func TestFibonacci(t *testing.T) {
	Convey("Testing fibonacci", t, func() {
		Convey("Fibo(0) should return 0", func() {
			So(fibo(0), ShouldEqual, 0)
		})
		Convey("Fibo(1) and fibo(2) should return 1", func() {
			So(fibo(1), ShouldEqual, 1)
			So(fibo(2), ShouldEqual, 1)
		})
		Convey("Fibo(8) should return 21", func() {
			So(fibo(8), ShouldEqual, 21)
		})
		Convey("Fibo(25) should return 75025", func() {
			So(fibo(25), ShouldEqual, 75025)
		})
	})
}
