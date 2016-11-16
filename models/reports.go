package models

import "log"

type Report struct{
	Name string `json:"name"`
	ID int64 `json:"id"`
	OwnerID int `json:"owner"`
	PlanID int `json:"flightplan"`
	Files []*FileMeta `json:"files"`
}

func (db *DB) GetReport(id int) (*Report, error) { 
	report := &Report{}
	err := db.QueryRow("SELECT * FROM reports WHERE id = ?", id).Scan(&report.ID, &report.Name, &report.OwnerID, &report.PlanID)
	if err != nil {
		return nil, err
	}
	report.Files, err = db.AllFilesForReport(report.ID)
	if err != nil {
		return nil, err
	}
	return report, nil
} 

func (db *DB) AllReportsForUser(userId int) ([]*Report, error) {
	rows, err := db.Query("SELECT * FROM reports WHERE owner = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	reports := make([]*Report, 0)
	for rows.Next() {
		report := new (Report)
		err := rows.Scan(&report.ID, &report.Name, &report.OwnerID, &report.PlanID)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

func (db *DB) AllReportsForPlan(planId int) ([]*Report, error) {
	rows, err := db.Query("SELECT * FROM reports WHERE plan = ?", planId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := make([]*Report, 0)
	for rows.Next() {
		report := new (Report)
		err := rows.Scan(&report.ID, &report.Name, &report.OwnerID, &report.PlanID)
		if err != nil {
			return nil, err
		}
		report.Files, err = db.AllFilesForReport(report.ID)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}
func (db *DB) AddReport(report Report) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO reports(name, owner, plan) VALUES (?,?,?)")
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(report.Name, report.OwnerID, report.PlanID)
	if err != nil {
		return -1, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	log.Printf("Added report: %s with ID %d. Rows affected: %d", report.Name, lastID, rowCnt)
	return lastID, nil
}
