package wrap

import "github.com/camsiabor/qcom/util"

type UtilWrap int;

const U UtilWrap = 1;

func (u UtilWrap) GetStr(o interface{}, def string, keys ... interface{}) string {
	return util.GetStr(o, def, keys...);
}




