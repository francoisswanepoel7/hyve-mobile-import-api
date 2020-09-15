package main

type CSVDataRaw struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	First_Name     string `json:"first_name"`
	Last_Name      string `json:"last_name"`
	Email          string `json:"email"`
	Tz             string `json:"tz"`
	Note           string `json:"note"`
	Datetime_Local string `json:"datetime_local"`
	Datetime_UTC   string `json:"datetime_utc"`
	Card           string `json:"card"`
}

type CSVDataValidated struct {
	CSVDataRaw
	Email_Valid bool   `json:"email_valid"`
	IP          string `json:"ip"`
	Processed   bool   `json:"processed"`
}

type OK struct {
	Result string `json:"result"`
}
