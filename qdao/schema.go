package qdao

import (
	"fmt"
	"github.com/camsiabor/qcom/util"
)

type VType string

const (
	VTYPE_BOOL   VType = "bool"
	VTYPE_INT    VType = "int"
	VTYPE_TIME   VType = "time"
	VTYPE_FLOAT  VType = "float"
	VTYPE_ARRAY  VType = "array"
	VTYPE_OBJECT VType = "object"
	VTYPE_STRING VType = "string"
	VTYPE_BINARY VType = "binary"
)

type VTypeCode int

const (
	VTYPE_CODE_BOOL   VTypeCode = 1
	VTYPE_CODE_INT    VTypeCode = 2
	VTYPE_CODE_TIME   VTypeCode = 3
	VTYPE_CODE_FLOAT  VTypeCode = 4
	VTYPE_CODE_ARRAY  VTypeCode = 5
	VTYPE_CODE_OBJECT VTypeCode = 6
	VTYPE_CODE_STRING VTypeCode = 7
	VTYPE_CODE_BINARY VTypeCode = 8
)

type VFactor int

const (
	VFACTOR_NORMAL VFactor = 0
	VFACTOR_KEY    VFactor = 1
)

type FieldSchema struct {
	Name     string
	Type     VType
	TypeCode VTypeCode
	Factor   VFactor
	Format   string
	Group    *GroupSchema
	Parent   *FieldSchema
	Sub      map[string]*FieldSchema
}

type GroupSchema struct {
	Name  string
	Key   *FieldSchema
	Field map[string]*FieldSchema
	DB    *DBSchema
}

type DBSchema struct {
	Name  string
	Group map[string]*GroupSchema
	Root  *Schema
}

type Schema struct {
	Name string
	DB   map[string]*DBSchema
}

func (field *FieldSchema) Init(name string, fieldOpts interface{}) {
	field.Name = name
	var vtype, ok = fieldOpts.(string)
	if ok {
		field.Type = VType(vtype)
		field.TypeCode = GetVTypeCode(field.Type)
		return
	}
	var mtype = util.AsMap(fieldOpts, false)
	if mtype == nil {
		panic(fmt.Errorf("unsupport group schema init %s | %v", name, fieldOpts))
	}
	vtype = util.GetStr(mtype, "", "vtype")
	if len(vtype) == 0 {
		vtype = string(VTYPE_STRING)
	}
	field.Type = VType(vtype)
	field.TypeCode = GetVTypeCode(field.Type)
	var bkey = util.GetBool(mtype, false, "key")
	if bkey {
		field.Factor = VFACTOR_KEY
	}
	field.Format = util.GetStr(mtype, "", "format")
	var subFieldOptions = util.GetMap(mtype, false, "schema")
	if subFieldOptions != nil {
		if field.Sub == nil {
			field.Sub = make(map[string]*FieldSchema)
		}
		for subFieldName, subFieldOpt := range subFieldOptions {
			if subFieldOpt == nil {
				continue
			}
			var sub = &FieldSchema{}
			sub.Init(subFieldName, subFieldOpt)
			sub.Parent = field
			field.Sub[subFieldName] = sub

		}
	}
}

func (group *GroupSchema) Init(name string, groupOptions map[string]interface{}) {
	group.Name = name
	if group.Field == nil {
		group.Field = make(map[string]*FieldSchema)
	}
	for fieldName, fieldOpts := range groupOptions {
		if fieldOpts == nil {
			continue
		}
		var field = &FieldSchema{}
		field.Init(fieldName, fieldOpts)
		if field.Factor == VFACTOR_KEY {
			if group.Key == nil {
				group.Key = field
			} else {
				panic(fmt.Errorf("key already defined %s => %s.%s", group.Key.Name, name, fieldName))
			}
		}
		field.Group = group
		group.Field[fieldName] = field
	}
}

func (db *DBSchema) Init(name string, dbOptions map[string]interface{}) {
	db.Name = name

	if db.Group == nil {
		db.Group = make(map[string]*GroupSchema)
	}

	for groupName, groupOpt := range dbOptions {
		if groupOpt == nil {
			continue
		}
		var groupMapOpt = util.AsMap(groupOpt, false)
		if groupMapOpt == nil {
			panic(fmt.Errorf("unsupport group opt %s = %v", groupName, groupMapOpt))
		}
		var group = &GroupSchema{}
		group.Init(groupName, groupMapOpt)
		group.DB = db
		db.Group[groupName] = group
	}
}

func (s *Schema) Init(name string, schemaOptions map[string]interface{}) {
	s.Name = name

	if s.DB == nil {
		s.DB = make(map[string]*DBSchema)
	}

	for dbName, dbOpt := range schemaOptions {
		if dbOpt == nil {
			continue
		}
		var dbMapOpt = util.AsMap(dbOpt, false)
		if dbMapOpt == nil {
			panic(fmt.Errorf("unsupport group opt %s = %v", dbName, dbMapOpt))
		}
		var db = &DBSchema{}
		db.Init(dbName, dbMapOpt)
		db.Root = s
		s.DB[dbName] = db
	}
}

func GetVTypeCode(vtype VType) VTypeCode {
	switch vtype {
	case VTYPE_BOOL:
		return VTYPE_CODE_BOOL
	case VTYPE_INT:
		return VTYPE_CODE_INT
	case VTYPE_TIME:
		return VTYPE_CODE_TIME
	case VTYPE_FLOAT:
		return VTYPE_CODE_FLOAT
	case VTYPE_ARRAY:
		return VTYPE_CODE_ARRAY
	case VTYPE_OBJECT:
		return VTYPE_CODE_OBJECT
	case VTYPE_STRING:
		return VTYPE_CODE_STRING
	case VTYPE_BINARY:
		return VTYPE_CODE_BINARY
	}
	panic(fmt.Errorf("unsupport vtype %s", vtype))
}
