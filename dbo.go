package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func initdb() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return db
}

func checkContactExists(id string) string {
	var contact_id = "-1"
	db := initdb()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	row := db.QueryRow("SELECT id FROM contact where id = " + id)
	q_err := row.Scan(&contact_id)
	if q_err != nil {

	}
	return contact_id
}

func checkContactTimeZoneExists(contact_id string) string {
	id := "-1"
	db := initdb()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	row := db.QueryRow("SELECT id FROM timezone where contact_id = " + contact_id)
	q_err := row.Scan(&id)
	if q_err != nil {

	}

	return id
}

func checkTimeZoneExists(tz string) string {
	var id string = "-1"

	db := initdb()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	row := db.QueryRow("SELECT id FROM timezone where tz = '" + tz + "'")
	q_err := row.Scan(&id)
	if q_err != nil {

	}

	return id

}

func insertTimeZone(tz string) string {
	db := initdb()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	country, region := getTZConstituents(tz)

	tzId := checkTimeZoneExists(tz)

	if tzId == "-1" {
		insert, err := db.Prepare("INSERT INTO timezone(country, region, tz) VALUES ( ?, ?, ?)")
		if err != nil {
			log.Fatalln(err)
		}

		insert.Exec(country, region, tz)
		if err != nil {
			log.Fatalln(err)
		}
		insert.Close()
	}

	tzId = checkTimeZoneExists(tz)
	return tzId

}

func insertCompleted(validated *CSVDataValidated) {
	db := initdb()
	defer db.Close()
	insert, err := db.Prepare("INSERT INTO processed(complete, contact_id) VALUES ( ?, ? )")
	if err != nil {
		log.Fatalln(err)
	}
	insert.Exec(1, validated.ID)

}

func insertExported(contact_id string) {
	db := initdb()
	defer db.Close()
	update, err := db.Prepare("UPDATE processed SET exported = ? WHERE contact_id = ?")
	if err != nil {
		log.Fatalln(err)
	}
	update.Exec(1, contact_id)

}

func getProcessed() []CSVDataValidated {
	db := initdb()
	defer db.Close()

	sql := "SELECT p.contact_id as id, c.title, c.first_name, c.last_name, c.email, c.note, c.card, c.ip, " +
		"t.tz, ct.tstamp_local, ct.tstamp_utc " +
		"FROM processed p " +
		"INNER JOIN contact c ON c.id = p.contact_id " +
		"INNER JOIN contact_timezone ct ON ct.contact_id = c.id " +
		"INNER JOIN timezone t ON t.id = ct.timezone_id " +
		"WHERE p.complete = 1 " +
		"ORDER BY p.contact_id ASC LIMIT 1000"

	results, err := db.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	var csv_export []CSVDataValidated
	for results.Next() {
		var export CSVDataValidated
		err = results.Scan(&export.ID, &export.Title, &export.First_Name, &export.Last_Name, &export.Email, &export.Note, &export.Card, &export.IP, &export.Tz, &export.Datetime_Local, &export.Datetime_UTC)
		if err != nil {
			panic(err.Error())
		}
		csv_export = append(csv_export, export)
		insertExported(export.ID)
	}
	return csv_export
}

func insertContact(validated *CSVDataValidated) bool {

	db := initdb()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	if checkContactExists(validated.ID) == "-1" {
		fmt.Println("importContact: " + validated.Email)
		insert, err := db.Prepare("INSERT INTO contact(id, title, first_name, last_name, email, note, card, ip) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )")
		if err != nil {
			log.Fatalln(err)
		}

		insert.Exec(validated.ID, validated.Title, validated.First_Name, validated.Last_Name, validated.Email, validated.Note, validated.Card, validated.IP)
		if err != nil {
			log.Fatalln(err)
		}
		insert.Close()
		return true
	}
	return false

}

func insertContactTimeZone(validated *CSVDataValidated, tzId string) {
	db := initdb()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	if checkContactTimeZoneExists(validated.ID) == "-1" {
		insert, err := db.Prepare("INSERT INTO contact_timezone(tstamp_local,tstamp_utc,contact_id, timezone_id) VALUES ( ?, ?, ?, ? )")
		if err != nil {
			log.Fatalln(err)
		}

		insert.Exec(validated.Datetime_Local, validated.Datetime_UTC, validated.ID, tzId)
		if err != nil {
			log.Fatalln(err)
		}
		insert.Close()
	}
}
