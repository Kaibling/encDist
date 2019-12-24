package libs

import (
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
    "time"
	"encoding/json"
)

// SQLiteInitPshDB ´sdd
func SQLiteInitPshDB(dbpath string) {
	 db, err := sql.Open("sqlite3", dbpath)
		CheckErr(err)		
        // insert
        stmt, err := db.Prepare("CREATE TABLE `publish` ( `uid` INTEGER PRIMARY KEY AUTOINCREMENT, `data` blob NULL,`created` DATE NULL, `identifier` VARCHAR(255) NULL );")
		 if err != nil {
			log.Println("Database already existing")
			db.Close()
			return 
		}
        _, err = stmt.Exec()
        if err != nil {
			log.Println(err)
		}
}


//SQLiteaddPublishData sds
func SQLiteaddPublishData(dbpath string, data CryptoData  ) string {

	jsonBlob, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("sqlite3", dbpath)
	CheckErr(err)
	guid := SHA1HashString(jsonBlob)
	// insert
    stmt, err := db.Prepare("INSERT INTO publish(data, created, identifier ) values(?,?,?)")
    CheckErr(err)
	res, err := stmt.Exec(jsonBlob, time.Now(),SHA1HashString(jsonBlob))
    CheckErr(err)
    id, err := res.LastInsertId()
	CheckErr(err)
	stmt.Close()
	log.Println(id)
	return guid
}


func SQLiteGetALLPublishData(dbpath string) {
	
	db, err := sql.Open("sqlite3", dbpath)
    CheckErr(err)
	// todo: has to be better
	rows, err := db.Query("Select uid,data,identifier from publish")
	CheckErr(err)
	var uid int
	var data []byte
	var identifier string
	var cnt int
	for rows.Next() {
		cnt++
		err = rows.Scan(&uid, &data,&identifier)
		CheckErr(err)
		log.Printf("%d %s",uid,identifier)
	}
	log.Println(cnt)
	rows.Close()
	db.Close()
}


func SQLiteGetPublishedData(dbpath string, hash string) []byte {
	
	db, err := sql.Open("sqlite3", dbpath)
    CheckErr(err)
    // todo: has to be better
    SQLQuery := "Select data from publish where identifier = '"+hash+"'"
    log.Debug(SQLQuery)
	rows, err := db.Query(SQLQuery)
	CheckErr(err)
	var data []byte
	for rows.Next() {
        err = rows.Scan(&data)
		CheckErr(err)
	}
	rows.Close()
    db.Close()
    return data
}

/*
func SQLiteBulkGetPublishedData(dbpath string, hash string) []byte {
}
*/