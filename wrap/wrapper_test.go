package wrap

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUtilWrap_GetStr(t *testing.T) {
	var v interface{};
	var refv = reflect.ValueOf(v);
	fmt.Println(refv.IsValid());
	//fmt.Println(refv.CanInterface());

	v = "";
	refv = reflect.ValueOf(v);
	fmt.Println(refv.IsValid());
	fmt.Println(refv.CanInterface());
}
