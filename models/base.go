package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

// CreateTable 创建表
func CreateTable(tables []interface{}, cover bool)  {
	for _, table := range tables {
		t := reflect.TypeOf(table)
		tableName := strings.Split(t.String(), ".")[1]
		if cover {
			if _, err := db.Exec("drop table if exists "+tableName+" cascade;"); err != nil {
				logrus.Warningln(err)
			}
		}
		sql := "CREATE TABLE " + tableName + " (\n"
		for i := 0; i < t.Elem().NumField(); i++ {
			field := t.Elem().Field(i)
			if i == t.Elem().NumField() - 1 {
				extra := field.Tag.Get("constraint")
				if extra != "" {
					sql += ",\n\t" + extra
				}
				break
			}
			if i != 0 {
				sql += ",\n\t"
			}
			sql += fmt.Sprintf(
				"%s %s %s",
				field.Tag.Get("json"),
				field.Tag.Get("type"),
				field.Tag.Get("constraint"),
			)
		}

		sql += "\n);"
		logrus.Debugln(sql)
		if _, err := db.Exec(sql); err != nil {
			logrus.Warningln(tableName, "创建失败:", err.Error())
		} else {
			logrus.Infoln(tableName, "创建成功")
		}

	}
}

// Insert 插入操作
func Insert(model interface{}) error {
	// 把存在的字段插入
	t := reflect.TypeOf(model)
	mp := make(map[string]string)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := reflect.ValueOf(model).Elem().FieldByName(field.Name).String()
		if value == "" {
			continue
		}
		mp[field.Tag.Get("json")] = value
	}
	tableName := strings.Split(t.String(), ".")[1]
	var kStr string
	var vStr string
	for k, v := range mp {
		kStr += k + ","
		vStr += "'" + v + "',"
	}
	kStr = kStr[:len(kStr)-1]
	vStr = vStr[:len(vStr)-1]
	sql := fmt.Sprintf("INSERT INTO %s(%s)\nVALUES (%s)", tableName, kStr, vStr)
	logrus.Debugln(sql)
	_, err := db.Exec(sql)
	return err
}

// getPrimaryKey 获取结构体主键名和对应值
func getPrimaryKey(model interface{}) (name string, value string, err error) {
	// 反射获取主键名和对应值
	t := reflect.TypeOf(model)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		constraint := field.Tag.Get("constraint")
		if strings.Contains(constraint, "PRIMARY KEY") {
			name = field.Tag.Get("json")
			value = reflect.ValueOf(model).Elem().FieldByName(field.Name).String()
			//value = reflect.ValueOf(model).Elem().FieldByName(name).String()
		}
	}
	if name != "" && value != "" {
		return
	}
	err = errors.New("PRIMARY KEY NO FOUND")
	return
}

// FindByKey 按主键查询表单
func FindByKey(model interface{}) error {
	name, value, err := getPrimaryKey(model)
	if err != nil {
		return err
	}
	// 构建查询语句并交付数据库查询
	odds :=  "WHERE " + name + " = '" + value + "'"
	return First(model, odds)
}

// First 按字段名查询单个
// odds => WHERE id = 1 || LIKE name = '%cj%' || ......
func First(model interface{}, odds string) error {
	// 反射获取表名
	t := reflect.TypeOf(model)
	tableName := strings.Split(t.String(), ".")[1]
	// 构建查询语句并交付数据库查询
	sql := "SELECT * FROM " + tableName + " " + odds
	return getFirst(model, sql)
}

// Find 按字段名查询多个
// odds => WHERE id = 2 || LIKE name = '%sh%' || ......
func Find(models interface{}, odds string) error {
	// 反射获取表名
	t := reflect.TypeOf(models).Elem()
	tableName := strings.Split(t.String(), ".")[1]
	// 构建查询语句并交付数据库查询
	sql := "SELECT * FROM " + tableName + " " + odds
	return getMany(models, sql)
}

func getFirst(model interface{}, sql string) error {
	// 1、将查询语句交付数据库查询
	logrus.Debugln(sql)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 2、开始解析返回的数据
	if !rows.Next() {
		return errors.New("sql: Scan called without calling Next")
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 3、构建指针数组，作为可变参数交付Scan函数进行赋值，获取数据
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(string)
	}
	if err := rows.Scan(values...); err != nil {
		return err
	}

	// 4、构建map，将数组数据转移到map缓存中
	m := make(map[string]interface{})
	for i, column := range columns {
		m[column] = *values[i].(*string)
	}

	// 5、使用反射将map映射到原结构体中
	t := reflect.TypeOf(model)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		tag := field.Tag.Get("json")
		if m[tag] != nil {
			v := reflect.ValueOf(model).Elem().FieldByName(field.Name)
			v.SetString(m[tag].(string))
		}
	}
	return nil
}

func getMany(models interface{}, sql string) error {
	// 1、将查询语句交付数据库查询
	logrus.Debugln(sql)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	var mList []map[string]interface{}
	for rows.Next() {
		// 2、开始解析返回的数据
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		// 3、构建指针数组，作为可变参数交付Scan函数进行赋值，获取数据
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(string)
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		// 4、构建map，将数组数据转移到map缓存中
		m := make(map[string]interface{})
		for i, column := range columns {
			m[column] = *values[i].(*string)
		}
		mList = append(mList, m)
	}

	// 反射获取表类型
	t :=  reflect.TypeOf(models).Elem().Elem()		// model.Table
	modelsValue := reflect.MakeSlice(reflect.TypeOf(models).Elem(), 0, len(mList))
	for _, m := range mList {
		newData := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("json")
			if m[tag] != nil {
				newData.FieldByName(field.Name).SetString(m[tag].(string))
			}
		}
		modelsValue = reflect.Append(modelsValue, newData)
	}
	v := reflect.ValueOf(models).Elem()
	v.Set(modelsValue)
	return nil
}

func Update(model interface{}) error {
	t := reflect.TypeOf(model)
	mp := make(map[string]string)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := reflect.ValueOf(model).Elem().FieldByName(field.Name).String()
		if value == "" {
			continue
		}
		mp[field.Tag.Get("json")] = value
	}
	name, value, err := getPrimaryKey(model)
	if err != nil {
		return err
	}
	tableName := strings.Split(t.String(), ".")[1]
	sql := "UPDATE " + tableName + " SET\n"
	tmp := 0
	for k, v := range mp {
		sql += fmt.Sprintf("\t%s = '%s'", k, v)
		tmp++
		if tmp != len(mp) {
			sql += ",\n"
		}
	}
	sql += fmt.Sprintf("\nWHERE %s = %s;", name, value)
	logrus.Debugln(sql)
	_, err = db.Exec(sql)
	return err
}

func Delete(model interface{}) error {
	t := reflect.TypeOf(model)
	tableName := strings.Split(t.String(), ".")[1]
	sql := "DELETE FROM " + tableName
	name, value, err := getPrimaryKey(model)
	if err != nil {
		return err
	}
	sql += fmt.Sprintf("\nWHERE %s = %s;", name, value)
	logrus.Debugln(sql)
	_, err = db.Exec(sql)
	return err
}
