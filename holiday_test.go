package main

import (
	"fmt"
	"testing"
)

func TestHolidayExhaustingInAMonth(t *testing.T) {

	balance := holidayBalanceByNow()
	fmt.Println(balance)
}
