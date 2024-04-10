package sqlutils

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func MysqlInitConfig(user, pass, address, port, dbname string) *mysql.Config {
	cfg := mysql.Config{
		User:      user,
		Passwd:    pass,
		Net:       "tcp",
		Addr:      strings.Join([]string{address, port}, ":"),
		DBName:    dbname,
		ParseTime: true,
	}
	return &cfg
}

func MysqlPing(cfg *mysql.Config) (*sql.DB, error) {
	var db *sql.DB

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

/* // Read action
func MyqlSelect(database *sql.DB, query string) ([]OcrRequest, error) {
	// An ocrRequests slice to hold data from returned rows.
	var ocrRequests []OcrRequest

	//query := "SELECT * FROM history WHERE doc_type=0 AND request_datetime BETWEEN '2024-02-01 00:00:00' AND '2024-02-20 23:59:59';"
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var request OcrRequest
		if err := rows.Scan(&request.Id, &request.Client_id, &request.Request_datetime, &request.Doc_type, &request.File_name, &request.Front_flag, &request.Result_code, &request.Ocr_result); err != nil {
			return nil, err
		}
		ocrRequests = append(ocrRequests, request)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ocrRequests, nil
} */

// Write action
func MysqlExec(database *sql.DB, query string) error {
	//doc_type := 0
	//query := fmt.Sprintf("UPDATE history SET request_datetime = DATE_ADD(request_datetime, INTERVAL 7 HOUR) WHERE doc_type=%d AND request_datetime BETWEEN '2024-02-01 00:00:00' AND '2024-02-20 23:59:59'", doc_type)

	fmt.Println(query)
	_, err := database.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
