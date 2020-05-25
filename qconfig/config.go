package qconfig

import (
	"encoding/json"
	"fmt"
	"github.com/camsiabor/qcom/util"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func ConfigLoad(filepath string, includename string, expand string) (config map[string]interface{}, err error) {

	var configfile *os.File
	configfile, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer configfile.Close()

	config = make(map[string]interface{})
	var decoder = json.NewDecoder(configfile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	if len(expand) > 0 {
		for key, val := range config {
			var valstr, ok = val.(string)
			if !ok {
				continue
			}
			var index = strings.Index(valstr, expand)
			if index != 0 {
				continue
			}
			var path = valstr[len(expand):]
			var subconfig, _ = ConfigLoad(path, includename, expand)
			if subconfig != nil {
				config[key] = subconfig
			}
		}
	}

	if len(includename) > 0 {
		var includes = util.GetMap(config, false, includename)
		if includes == nil {
			return config, err
		}
		for key, val := range includes {
			if val == nil {
				continue
			}
			var sval, ok = val.(string)
			if ok {
				subconfig, _ := ConfigLoad(sval, includename, expand)
				if subconfig != nil {
					config[key] = subconfig
				}
			}
		}
	}

	return config, err
}

func ConfigParse(path string) (config map[string]interface{}, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// Useful for command line to override options specified in config file
// Debug is not updated.
func ConfigUpdate(old, new map[string]interface{}) {
	// Using reflection here is not necessary, but it's a good exercise.
	// For more information on reflections in Go, read "The Laws of Reflection"
	// http://golang.org/doc/articles/laws_of_reflection.html
	newVal := reflect.ValueOf(new).Elem()
	oldVal := reflect.ValueOf(old).Elem()

	// typeOfT := newVal.Type()
	for i := 0; i < newVal.NumField(); i++ {
		newField := newVal.Field(i)
		oldField := oldVal.Field(i)
		// log.Printf("%d: %s %s = %v\n", i,
		// typeOfT.Field(i).Name, newField.Type(), newField.Interface())
		switch newField.Kind() {
		case reflect.Interface:
			if fmt.Sprintf("%v", newField.Interface()) != "" {
				oldField.Set(newField)
			}
		case reflect.String:
			s := newField.String()
			if s != "" {
				oldField.SetString(s)
			}
		case reflect.Int:
			i := newField.Int()
			if i != 0 {
				oldField.SetInt(i)
			}
		}
	}

}
