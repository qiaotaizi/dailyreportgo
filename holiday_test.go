package main

import (
	"fmt"
	"testing"
)

func TestHolidayExhaustingInAMonth(t *testing.T){
	fmt.Println(exhausting)
	fmt.Println(exhausted)
	//ti,_:=time.Parse("2006-01-02","2019-12-18")
	balance:=holidayBalanceByNow()
	fmt.Println(balance)
}
