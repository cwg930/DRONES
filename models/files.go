package models

import "log"

type FileMeta struct{
	FileName string `json:"filename"`
	ID int `json:"id"`
	OwnerID int `json:"owner"`
	ReportID int `json:"report"`
	PointID int `json:"point"`
}

func (db *DB) AllFiles(ownerID int) ([]*FileMeta, error) {
	rows, err := db.Query("SELECT * FROM files WHERE owner = ?", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*FileMeta, 0)
	for rows.Next() {
		file := new (FileMeta)
		err := rows.Scan(&file.ID, &file.OwnerID, &file.ReportID, &file.PointID, &file.FileName)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return files,  nil
}

func (db *DB) AllFilesForReport(reportID int) ([]*FileMeta, error) {
	rows, err := db.Query("SELECT * FROM files WHERE report = ?", reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]*FileMeta, 0)
	for rows.Next() {
		file := new (FileMeta)
		err := rows.Scan(&file.ID, &file.OwnerID, &file.ReportID, &file.PointID, &file.FileName)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}


func (db *DB) GetFile(id int) (*FileMeta, error) {
	file := &FileMeta{}
	err := db.QueryRow("SELECT * FROM files WHERE id = ?", id).Scan(&file.ID, &file.OwnerID, &file.ReportID, &file.PointID, &file.FileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}
//add a reference to a file 
func (db *DB) AddFile(file FileMeta) error {
	stmt, err := db.Prepare("INSERT INTO files(owner,report,point,filename) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(file.OwnerID, file.ReportID, file.PointID, file.FileName)
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
	log.Printf("Added file: %s with ID %d. Rows affected: %d", file.FileName, lastID, rowCnt)
	return nil
}
