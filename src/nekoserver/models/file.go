package models

import (
	"fmt"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

func FetchFileList() (error, []data.NekohandFile) {
	statement := fmt.Sprintf("select * from files")
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err, nil
	}
	rows, err := db.Query(statement)

	if err != nil {
		return err, []data.NekohandFile{}
	}

	nfiles := []data.NekohandFile{}

	for rows.Next() {
		var nf data.NekohandFile
		if err = rows.Scan(&nf.FID, &nf.FileId, &nf.HashId, &nf.FileName, &nf.CreatedAt, &nf.ModifiedAt); err != nil {
			return err, nil
		}
		nfiles = append(nfiles, nf)
	}
	return nil, nfiles
}
