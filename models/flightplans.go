package models

import "log"

type FlightPlan struct{
	Name string `json:"name"`
	ID int `json:"id"`
	OwnerID int `json:"owner"`
	Points []*Point `json:"points"`
}

type Point struct{
	ID int `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt"`
	Pic bool `json:"pic"`
}

func (db *DB) GetPlan(id int) (*FlightPlan, error) {
	plan := &FlightPlan{}
	err := db.QueryRow("SELECT * FROM flightplans WHERE id = ?", id).Scan(&plan.ID, &plan.Name, &plan.OwnerID)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id,lat,lon,alt,picture FROM points WHERE plan = ?", plan.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := make([]*Point, 0)
	for rows.Next() {
		point := new(Point)
		err := rows.Scan(&point.ID, &point.Lat, &point.Lon, &point.Alt, &point.Pic)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	} 
	plan.Points = points
	return plan, nil
}

func (db *DB) AllPlansForUser(ownerID int) ([]*FlightPlan, error){
	rows, err := db.Query("SELECT id, name FROM flightplans WHERE owner = ?", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	plans := make([]*FlightPlan,0)
	for rows.Next() {
		plan := new(FlightPlan)
		err := rows.Scan(&plan.ID, &plan.Name)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return plans, nil
}

func (db *DB) AddFlightPlan(plan FlightPlan) error {
	stmt, err := db.Prepare("INSERT INTO flightplans(name,owner) VALUES (?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(plan.Name, plan.OwnerID)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Added flightplan: %s with ID %d. Rows affected: %d", plan.Name, lastID, rowCnt)
	return nil
}

func (db *DB) AddAllPoints(planID int, points []*Point) error {
	sqlStr := "INSERT INTO points(plan, lat, lon, alt, picture) VALUES "
	vals := []interface{}{}
	
	for _, row := range points {
		sqlStr += "(?,?,?,?,?),"
		vals = append(vals, planID, row.Lat, row.Lon, row.Alt, row.Pic)
	}
	sqlStr = sqlStr[0:len(sqlStr)-2]
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(vals)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Added points, lastID=%d. Rows affected: %d", lastID, rowCnt)
	return nil
}
