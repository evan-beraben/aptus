package data_access

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var _db *sql.DB
var _err error
var _host string
var _dbName string
var _uname string
var _pw string

func Init(dbUrl string, uname string, pw string) error{
	_host, _dbName = ParseUrl(dbUrl)
	_uname = uname
	_pw = pw

	return Connect()
}

func Connect() error {
	log.Printf("Connecting to mysql %s", _host)
	_db, _err = sql.Open("mysql", _uname+":"+_pw+"@tcp("+_host+")"+"/"+_dbName)

	if _err != nil {
		log.Println("Failed to connect to db")
		return _err
	}

	log.Println("Mysql connection success.")
	return nil
}

func Close() {
	if _db != nil {
		log.Println("Closing mysql connection")
		_db.Close()
	}
}

func ParseUrl(dbUrl string) (string, string) {
	split := strings.Split(dbUrl, "?")
	if len(split) == 0 {
		log.Fatal("DB config parse error")
	}

	split1 := strings.Split(split[0], "jdbc:mysql://")
	if len(split1) != 2 {
		log.Fatal("DB config parse error")
	}

	split3 := strings.Split(split1[1], "/")
	if len(split3) != 2 {
		log.Fatal("DB config parse error")
	}

	return split3[0], split3[1]

}

func ExecSQL(sql string) error{
	r, err := _db.Exec(sql)
	if err != nil{
		log.Printf("Failed to execute sql:\n%s\n%s", sql, err)
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil{
		log.Printf("Failed to get result for sql:\n%s\n%s", sql, err)
		return err
	}

	log.Printf("%s\nq%d Rows affected\n", sql, rowsAffected)

	return err;
}

func Select(t reflect.Type, from string, where string) (interface{}, error) {

	sel := buildSelect(t, "a")
	sel += " FROM " + from
	if where != ""{
		sel += " WHERE " + where
	}


	//log.Println("Running query: " + sel)
	rows, err := _db.Query(sel)

	if err != nil {
		log.Println("Failed to query rows: " + err.Error())
		return nil, err
	}

	slice := make([]interface{}, 0, 0)

	if rows == nil {
		return slice, nil
	}

	cols, err := rows.Columns()
	vals := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(vals))
	for c, _ := range cols {
		scanArgs[c] = &vals[c]
	}

	for rows.Next() {

		item := reflect.New(t).Interface()
		itemVal := reflect.ValueOf(item)
		itemElem := itemVal.Elem()

		rows.Scan(scanArgs...)

		scanResult(itemElem, vals, new(int))

		slice = append(slice, item)

	}

	return slice, nil
}


func scanResult(itemElem reflect.Value, vals []sql.RawBytes, valPos *int){

	for i := 0; i < itemElem.NumField(); i++ {
		f := itemElem.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			scanResult(f, vals, valPos)
			break
		case reflect.String:
			s := string(vals[*valPos])
			f.SetString(s)
			*valPos++
			break
		case reflect.Int:
			s := string(vals[*valPos])
			vi, err :=strconv.ParseInt(s,10,32)
			if err != nil {
				log.Fatal("failed to convert Int")
			}
			f.SetInt(vi)
			*valPos++
			break
		default:
			break
		}

	}
}

func buildSelect(t reflect.Type, alias string) string {
	result := "SELECT "

	fields := buildMap(t, alias)

	for i := 0; i < len(fields); i++ {
		if i > 0 {
			result += ","
		}
		result += fields[i]
	}

	return result;
}

func buildMap(t reflect.Type, prefix string) []string {
	result := []string{}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {

			r := buildMap(t.Field(i).Type, t.Field(i).Name)
			//append result
			for a := 0; a < len(r); a++ {
				result = append(result, r[a])
			}
		} else {
			result = append(result, prefix+"."+t.Field(i).Name)
		}
	}

	return result
}

func getSelectResult() {

}

//func Select(fields interface{}, from string, where string){
/*
	query := "SELECT "
	for i, s := range fields{
		query += s
		if i < (len(fields) - 1){
			query += ","
		}
		query += " "
	}

	query += "FROM " + from;
	if len(where) > 0 {
		query += "WHERE " + where
	}

	rows, err := _db.Query(query)

	if(err != nil){
		log.Println("Query failed")
		return
	}

	if rows != nil{
		rows.sca
	}
*/

//}
