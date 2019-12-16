package main


import (
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"time"
	"crypto/rsa"
	"encoding/json"
)

// SQLiteInitDB ´sdd
func SQLiteInitDB(dbpath string) {
	 db, err := sql.Open("sqlite3", dbpath)
		checkErr(err)		
        // insert
        stmt, err := db.Prepare("CREATE TABLE `user` ( `uid` INTEGER PRIMARY KEY AUTOINCREMENT,`username` VARCHAR(64) NULL, `privateKey` blob NULL,`created` DATE NULL,`password` VARCHAR(64) NULL);")
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

//SQLiteaddUser sds
func SQLiteaddUser(dbpath string, user *User, privKey *rsa.PrivateKey) {

	jsonPrivKey, err := json.Marshal(privKey)
	if err != nil {
		log.Println(err)
	}

	//encrypt privateKey
	encyptPrivKey := AESencryptData([]byte(user.password),jsonPrivKey)
	 db, err := sql.Open("sqlite3", dbpath)

    checkErr(err)
    // insert
    stmt, err := db.Prepare("INSERT INTO user(username, privateKey, created,password ) values(?,?,?,?)")
    checkErr(err)
	res, err := stmt.Exec(user.name, encyptPrivKey, time.Now(), user.password)
    checkErr(err)
    id, err := res.LastInsertId()
	checkErr(err)
	stmt.Close()
    log.Println(id)
}


func SQLiteGetUser(dbpath string, name string) *User {
	
	db, err := sql.Open("sqlite3", dbpath)
    checkErr(err)

	// todo: has to be better
	rows, err := db.Query("Select username,password from user where username = '" + name + "'")
	checkErr(err)
	var username string
    var password string
	var cnt int
	for rows.Next() {
		cnt++
		err = rows.Scan(&username, &password)
        checkErr(err)
	}
	rows.Close()
	db.Close()
	
	if cnt == 0 {
		return nil
	}
	returnUser := new(User)
	returnUser.name = username
	returnUser.password = password
	return returnUser
}

func SQLiteGetKeyPair(dbpath string, name string, aeskey string) (*rsa.PrivateKey,error) {
	
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	
    var encPrivateKeyJSON []byte
	rows, err := db.Query("Select privateKey from user where username = '" + name + "'")
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&encPrivateKeyJSON)
    	checkErr(err)
	}
	rows.Close()
	db.Close()

	privateKeyJSON := AESdecryptdata([]byte(aeskey),encPrivateKeyJSON)
	var ha rsa.PrivateKey

	err = json.Unmarshal(privateKeyJSON,&ha)
	if err != nil {
		log.Errorln(err)
	}
	return &ha, nil
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}