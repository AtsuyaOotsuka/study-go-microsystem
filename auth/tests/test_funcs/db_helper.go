package test_funcs

import (
	"database/sql"
	"microservices/auth/tests/test_db_seeder"

	"golang.org/x/crypto/bcrypt"
)

func truncateTable(db *sql.DB, tableName string) error {
	query := "TRUNCATE TABLE " + tableName
	_, err := db.Exec(query)
	return err
}

func DbCleanup(db *sql.DB) ([]DbRecords, error) {
	truncateTable(db, "users")
	dbRecords, err := CreateSeeders(db)
	if err != nil {
		return nil, err
	}
	return dbRecords, nil
}

type DbRecords struct {
	TableName string
	Count     int
	Data      []map[string]interface{}
}

func CreateSeeders(db *sql.DB) ([]DbRecords, error) {
	dbRecords := []DbRecords{}

	users := test_db_seeder.GetUsersSeeders(5, false)
	for _, user := range users {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		Insert, err := db.Exec("INSERT INTO users (name, email, password, refresh_token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			user.Name, user.Email, password, user.RefreshToken, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		InsertId, err := Insert.LastInsertId()
		if err != nil {
			return nil, err
		}
		dbRecords = append(dbRecords, DbRecords{
			TableName: "users",
			Count:     len(users),
			Data: []map[string]interface{}{
				{
					"id":            InsertId,
					"name":          user.Name,
					"email":         user.Email,
					"password":      user.Password, // 平文のまま保存
					"refresh_token": user.RefreshToken,
					"created_at":    user.CreatedAt,
					"updated_at":    user.UpdatedAt,
					"deleted_at":    user.DeletedAt,
				},
			},
		})
	}
	return dbRecords, nil
}

func ExistsRecord(db *sql.DB, table string, filter map[string]interface{}) bool {
	record := GetRecords(db, table, filter)
	return len(record) > 0
}

func GetRecords(db *sql.DB, table string, filter map[string]interface{}) []DbRecords {
	records := []DbRecords{}
	query := "SELECT * FROM " + table + " WHERE "
	args := []interface{}{}
	i := 0
	for k, v := range filter {
		if i > 0 {
			query += " AND "
		}
		query += k + " = ?"
		args = append(args, v)
		i++
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return records
	}
	cols, err := rows.Columns()
	if err != nil {
		return records
	}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return records
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		records = append(records, DbRecords{
			TableName: table,
			Data:      []map[string]interface{}{m},
		})
	}
	return records
}
