package core

import (
	"encoding/json"
	"strconv"
	"strings"
)

type UserStat struct {
	Name  string
	Value int64
}
type KeyValue struct {
	Name  string
	Value string
}
type ComplexStat struct {
	Name   string
	Values []KeyValue
}

type FieldStruct struct {
	Name string
	Type string
	Size int
}
type ComplexStatTableInfo struct {
	Name      string
	Fields    []FieldStruct
	AutoCount bool
}

func GetComplexStatsTableInfos(table string) (ComplexStatTableInfo, error) {
	var values string
	var info ComplexStatTableInfo
	err := Db.QueryRow("SELECT tableName, fields, autoCount FROM neb_stats_tables WHERE tableName = ?", table).Scan(&info.Name, &values, &info.AutoCount)
	if err != nil {
		Error.Println("Could not read ComplexStatsTableInfos: ", err)
		return info, err
	}
	json.Unmarshal(([]byte)(values), &info.Fields)

	return info, nil
}

func GetComplexStatsTablesInfos() ([]ComplexStatTableInfo, error) {
	var ret = make([]ComplexStatTableInfo, 0)

	rows, err := Db.Query("SELECT tableName, fields, autoCount FROM neb_stats_tables")
	defer rows.Close()

	if err != nil {
		return ret, err
	}

	for rows.Next() {
		var info ComplexStatTableInfo
		var fields string
		err = rows.Scan(&info.Name, &fields, &info.AutoCount)
		if err != nil {
			return ret, err
		}
		json.Unmarshal(([]byte)(fields), &info.Fields)

		ret = append(ret, info)
	}

	err = rows.Err()

	if err != nil {
		return ret, err
	}
	return ret, nil
}

func GetUserStatsFields() ([]string, error) {
	var statFields []string

	tables, err := GetComplexStatsTablesInfos()
	if err != nil {
		Error.Println("Could not read ComplexStatsTableInfos: ", err)
		return statFields, err
	}

	for _, table := range tables {
		if table.Name == "users" {
			for _, field := range table.Fields {
				statFields = append(statFields, field.Name)
			}
		} else if table.AutoCount {
			statFields = append(statFields, table.Name)
		}
	}

	return statFields, nil
}
func SetUsersStatFields(fields string) error {
	var utable ComplexStatTableInfo
	for _, field := range strings.Split(fields, ",") {
		var tfield FieldStruct
		tfield.Name = field
		tfield.Size = 11
		tfield.Type = "int"
		utable.Fields = append(utable.Fields, tfield)
	}

	res, err := json.Marshal(utable.Fields)
	if err != nil {
		return err
	}

	_, err = Db.Exec("UPDATE neb_stats_tables SET fields=? WHERE tableName = 'users'", res)

	return err
}

func SetStatFields(table ComplexStatTableInfo) error {
	res, err := json.Marshal(table.Fields)
	if err != nil {
		return err
	}

	_, err = Db.Exec("UPDATE neb_stats_tables SET fields=?, autoCount=? WHERE tableName = ?", res, table.AutoCount, table.Name)

	return err
}

func AddUserStatFields(name string) error {
	var fields string
	err := Db.QueryRow("SELECT fields FROM neb_stats_table WHERE tableName = users").Scan(&fields)
	if err != nil {
		Error.Println("Could not get fields from neb_stats_table for users")
		return err
	}

	tFields := strings.Split(fields, ",")
	for _, field := range tFields {
		if field == name {
			return &NebuleuseError{Code: NebError, Msg: "Field already exists in users"}
		}
	}
	fields += "," + name
	_, err = Db.Exec("UPDATE fields FROM neb_stats_table WHERE tableName = users", fields)
	if err != nil {
		return &NebuleuseError{Code: NebError, Msg: "Could not update users stats fields"}
	}

	return nil
}

func DeleteUserStatFields(name string) error {
	var fields string
	err := Db.QueryRow("SELECT fields FROM neb_stats_table WHERE tableName = users").Scan(&fields)
	if err != nil {
		Error.Println("Could not get fields from neb_stats_table for users")
		return err
	}

	var newFields string
	tFields := strings.Split(fields, ",")
	for _, field := range tFields {
		if field != name {
			newFields += field + ","
		}
	}

	//Remove trailing ,
	newFields = newFields[:len(newFields)]

	_, err = Db.Exec("UPDATE fields FROM neb_stats_table WHERE tableName = users", newFields)
	if err != nil {
		return &NebuleuseError{Code: NebError, Msg: "Could not update users stats stats"}
	}

	return nil
}
func getTypeForQuery(Type string, size int) string {
	query := ""
	switch Type {
	case "string":
		query += "varchar("
		query += strconv.Itoa(size)
		query += ")"
	case "int":
		query += "int("
		query += strconv.Itoa(size)
		query += ")"
	case "text":
		query += "text"
	case "timestamp":
		query += "timestamp"
	default:
		query += "int("
		query += strconv.Itoa(size)
		query += ")"
	}
	return query
}
func AddStatTable(table ComplexStatTableInfo) error {
	query := "CREATE TABLE neb_users_stats_"
	query += table.Name
	query += " ( "
	for _, field := range table.Fields {
		query += field.Name
		query += " "
		query += getTypeForQuery(field.Type, field.Size)
		query += ","
	}
	query = query[:len(query)-1]
	query += " );"

	_, err := Db.Exec(query)
	if err != nil {
		return err
	}

	res, err := json.Marshal(table.Fields)
	if err != nil {
		return err
	}
	_, err = Db.Exec("INSERT INTO neb_stats_tables VALUES(?,?,?)", table.Name, res, table.AutoCount)
	if err != nil {
		return err
	}

	return nil
}

func SetStatTable(table ComplexStatTableInfo) error {
	oldTable, err := GetComplexStatsTableInfos(table.Name)
	if err != nil {
		return err
	}

	removeQuery := "ALTER TABLE neb_users_stats_" + table.Name + " DROP COLUMN "
	updateQuery := "ALTER TABLE neb_users_stats_" + table.Name + " CHANGE "
	if err != nil {
		return err
	}
	for _, oldField := range oldTable.Fields {
		found := false
		update := false
		var newType FieldStruct
		for _, newField := range table.Fields {
			if newField.Name == oldField.Name {
				found = true
				if newField.Type != oldField.Type || newField.Size != oldField.Size {
					update = true
					newType = newField
				}
			}
		}
		if !found {
			_, err := Db.Exec(removeQuery + oldField.Name)
			if err != nil {
				return err
			}
		} else if update {
			_, err := Db.Exec(updateQuery + newType.Name + " " + newType.Name + " " + getTypeForQuery(newType.Type, newType.Size))
			if err != nil {
				return err
			}
		}
	}
	addQuery := "ALTER TABLE neb_users_stats_" + table.Name + " ADD COLUMN "
	for _, newField := range table.Fields {
		found := false
		for _, oldField := range oldTable.Fields {
			if oldField.Name == newField.Name {
				found = true
			}
		}
		if !found {
			_, err := Db.Exec(addQuery + newField.Name + " " + getTypeForQuery(newField.Type, newField.Size))
			if err != nil {
				return err
			}
		}
	}
	err = SetStatFields(table)

	return err
}

func DeleteStatTable(name string) error {
	query := "DROP TABLE neb_users_stats_"
	query += name

	_, err := Db.Exec(query)
	if err != nil {
		return err
	}

	_, err = Db.Exec("DELETE FROM neb_stats_tables WHERE tableName = ?", name)
	if err != nil {
		return err
	}

	return nil
}
