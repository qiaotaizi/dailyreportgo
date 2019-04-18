package holiday

import (
	"fmt"
	"testing"
	"time"
)

func TestHolidayExhaustingInAMonth(t *testing.T){
	ti,_:=time.Parse("2006-01-02","2019-12-18")
	b:=HolidayExhaustingInAMonth(ti)
	fmt.Println(b)
}
