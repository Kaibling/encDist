package libs

import (
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
    "time"
    	"crypto/rsa"
	"encoding/json"
)
// SQLiteInitTknDB ´sdd
func SQLiteInitTknDB(dbpath string) {
	 db, err := sql.Open("sqlite3", dbpath)
		CheckErr(err)		
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
	encyptPrivKey, err := AESencryptData([]byte(user.Password),jsonPrivKey)
	if err != nil {
		log.Debug(err)
		return
	}
	 db, err := sql.Open("sqlite3", dbpath)

    CheckErr(err)
	// insert
	passwordHash := SHA1HashString([]byte(user.Password))
    stmt, err := db.Prepare("INSERT INTO user(username, privateKey, created,password ) values(?,?,?,?)")
    CheckErr(err)
	res, err := stmt.Exec(user.Name, encyptPrivKey, time.Now(), passwordHash)
    CheckErr(err)
    id, err := res.LastInsertId()
	CheckErr(err)
	stmt.Close()
    log.Println(id)
}


func SQLiteGetFullUser(dbpath string, name string,password string) *FullUser {
	
	db, err := sql.Open("sqlite3", dbpath)
    CheckErr(err)

	passwordHash := SHA1HashString([]byte(password))
	// todo: has to be better
	rows, err := db.Query("Select username,password,privateKey from user where username = '" + name + "' and password = '"+ passwordHash+"'")
	CheckErr(err)
	var username string
	var pwd string
	var encPrivateKeyJSON []byte
	var cnt int
	for rows.Next() {
		cnt++
		err = rows.Scan(&username, &pwd,&encPrivateKeyJSON)
        CheckErr(err)
	}
	rows.Close()
	db.Close()
	
	if cnt == 0 {
		return nil
	}

	privateKeyJSON,err := AESdecryptdata([]byte(password),encPrivateKeyJSON)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	var ha rsa.PrivateKey

	err = json.Unmarshal(privateKeyJSON,&ha)
	if err != nil {
		log.Errorln(err)
	}
	returnUser := new(FullUser)
	returnUser.Name = username
	returnUser.Password = password
	returnUser.PrivateKey = &ha
	return returnUser
}

func SQLiteGetUser(dbpath string, name string) *User {
	
	db, err := sql.Open("sqlite3", dbpath)
    CheckErr(err)

	// todo: has to be better
	rows, err := db.Query("Select username,password from user where username = '" + name + "'")
	CheckErr(err)
	var username string
    var password string
	var cnt int
	for rows.Next() {
		cnt++
		err = rows.Scan(&username, &password)
        CheckErr(err)
	}
	rows.Close()
	db.Close()
	
	if cnt == 0 {
		return nil
	}
	
	returnUser := new(User)
	returnUser.Name = username
	returnUser.Password = password
	return returnUser
}

func SQLiteGetKeyPair(dbpath string, name string, aeskey string) (*rsa.PrivateKey,error) {
	
	db, err := sql.Open("sqlite3", dbpath)
	CheckErr(err)
	
    var encPrivateKeyJSON []byte
	rows, err := db.Query("Select privateKey from user where username = '" + name + "'")
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&encPrivateKeyJSON)
    	CheckErr(err)
	}
	rows.Close()
	db.Close()

	privateKeyJSON,err := AESdecryptdata([]byte(aeskey),encPrivateKeyJSON)
	if err != nil {
		log.Errorln(err)
		return nil,nil
	}
	var ha rsa.PrivateKey

	err = json.Unmarshal(privateKeyJSON,&ha)
	if err != nil {
		log.Errorln(err)
	}
	return &ha, nil
}
