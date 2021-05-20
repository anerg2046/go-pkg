package utils

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type JsonArray []interface{}

// type JsonArray json.RawMessage

func (jsonArr *JsonArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, &jsonArr)
	return err
}

func (jsonArr JsonArray) Value() (driver.Value, error) {
	if jsonArr == nil {
		return "[]", nil
	} else {
		d, _ := json.Marshal(jsonArr)
		return string(d), nil
	}
}

func JsonArrayFormat(value string) JsonArray {
	var _jsonArr JsonArray
	value = strings.TrimSpace(value)
	if value == "" {
		return _jsonArr
	}
	arr := strings.Split(value, ",")
	for _, i := range arr {
		i := strings.TrimSpace(i)
		r, _ := regexp.Compile(`[0-9]+`)
		if ok := r.MatchString(i); ok {
			in, _ := strconv.Atoi(i)
			_jsonArr = append(_jsonArr, in)
		} else {
			_jsonArr = append(_jsonArr, i)
		}
	}
	return _jsonArr
}

func (JsonArray) GormDataType() string {
	return "json"
}

func (JsonArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (jsonArr JsonArray) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	jsonStr, _ := json.Marshal(jsonArr)
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{jsonStr},
	}
}

type JsonArrayQueryExpression struct {
	column   string
	hasValue bool
	values   []interface{}
}

func JsonArrayQuery(column string) *JsonArrayQueryExpression {
	return &JsonArrayQueryExpression{column: column}
}

func (jsonQuery *JsonArrayQueryExpression) HasValue(value ...interface{}) *JsonArrayQueryExpression {
	jsonQuery.values = value
	jsonQuery.hasValue = true
	return jsonQuery
}

func (jsonQuery *JsonArrayQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "mysql", "sqlite":
			switch {
			case jsonQuery.hasValue:
				if len(jsonQuery.values) > 0 {
					values, _ := json.Marshal(jsonQuery.values)
					builder.WriteString(fmt.Sprintf("JSON_CONTAINS(%s, '%s')", stmt.Quote(jsonQuery.column), values))
				}
			}
		}
	}
}
