package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type Employees []model.Employee

func main() {
	if err := godotenv.Load(); err != nil {
		errorExit("failed to load .env file: %v", err)
	}

	db, err := sqlx.Connect("mysql", os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		errorExit("failed to connect to database: %v", err)
	}
	defer db.Close()

	f, err := os.Open("employees.json")
	if err != nil {
		errorExit("failed to open employees.json: %v", err)
	}
	defer f.Close()

	var employees Employees
	if err := json.NewDecoder(f).Decode(&employees); err != nil {
		errorExit("failed to decode employees.json: %v", err)
	}

	for _, emp := range employees {
		_, err := db.Exec(`
			INSERT INTO employees (id, name, position, branch_office, employment_type, is_excluded, guaranteed_doorprize, present_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, emp.ID, emp.NamaKaryawan, emp.Jabatan, emp.KantorCabang, emp.JenisKepegawaian, emp.IsExcluded == 1, emp.GuaranteedDoorprize == 1, nil)
		if err != nil {
			errorExit("failed to insert employee %s: %v", emp.NamaKaryawan, err)
		}
	}

	log.Println("seeding completed successfully")
}

func errorExit(msg string, args ...any) {
	log.Printf(msg+"\n", args...)
	os.Exit(1)
}
