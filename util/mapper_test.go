package util

import (
	"fmt"
	"testing"
)

func TestMapper(t *testing.T) {

	var mapperConfig, err = ConfigLoad("mapping.json", "");
	if (err != nil) {
		panic(err);
	}


	var manager = GetMapperManager();
	manager.Init(mapperConfig);

	var mapper = manager.Get("tushare.khistory");

	var m = make(map[string]interface{});
	m["ts_code"] = "000001.SZ";

	fmt.Println(m);
	mapper.Map(m, false);
	fmt.Println(m);


}