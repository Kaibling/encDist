package libs


import (
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"time"
	"crypto/rsa"
	"encoding/json"
	"os"
		"flag"
		//"encoding/hex"
)

type Configuration struct {
	BindingIPAddress string
	BindingPort      string
	DBpath	         string
}

type PlainDataTransfer struct {
	Token string
	Data []byte
}

type CryptoDataTransfer struct {
	Token string
	CryptoData CryptoData
}

func ParseArguments() map[string]string {

	arguments := make(map[string]string)
	configString := flag.String("config", "config.json", "filepath to configuration file")
	flag.Parse()

	arguments["configFilePath"] = *configString
	log.Print("load command line arguments ")
	log.Print(arguments)

	return arguments
}

func ParseConfigurationFile(filepath string) *Configuration {

	if filepath == "config.json" {
		//default path found
		//create config file
		_, err := os.Stat(filepath)

		if os.IsNotExist(err) {
			log.Println("file does not exists")

			fo, err := os.Create("config.json")
			checkErr(err)

			returnConfig := new(Configuration)
			returnConfig.BindingIPAddress = "127.0.0.1"
			returnConfig.BindingPort = "8070"
			returnConfig.DBpath = "./sqlitehier.db"

			configString, err := json.Marshal(returnConfig)
			checkErr(err)

			_, err = fo.Write(configString)
			checkErr(err)

			defer fo.Close()
			log.Println("Configuration file created")
			return returnConfig
		}

	}
	log.Println("opening configuration file: " + filepath)
	returnConfig := new(Configuration)
	file, err := os.Open(filepath)
	checkErr(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&returnConfig)
	checkErr(err)

	return returnConfig
}

//User sds
type FullUser struct {
	Name string
	Password string
	PrivateKey *rsa.PrivateKey
}

//User sds
type User struct {
	Name string
	id int32
	Password string
}


// SQLiteInitTknDB ´sdd
func SQLiteInitTknDB(dbpath string) {
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
	encyptPrivKey := AESencryptData([]byte(user.Password),jsonPrivKey)
	 db, err := sql.Open("sqlite3", dbpath)

    checkErr(err)
	// insert
	passwordHash := SHA1HashString([]byte(user.Password))
    stmt, err := db.Prepare("INSERT INTO user(username, privateKey, created,password ) values(?,?,?,?)")
    checkErr(err)
	res, err := stmt.Exec(user.Name, encyptPrivKey, time.Now(), passwordHash)
    checkErr(err)
    id, err := res.LastInsertId()
	checkErr(err)
	stmt.Close()
    log.Println(id)
}


func SQLiteGetFullUser(dbpath string, name string,password string) *FullUser {
	
	db, err := sql.Open("sqlite3", dbpath)
    checkErr(err)

	passwordHash := SHA1HashString([]byte(password))
	// todo: has to be better
	rows, err := db.Query("Select username,password,privateKey from user where username = '" + name + "' and password = '"+ passwordHash+"'")
	checkErr(err)
	var username string
	var pwd string
	var encPrivateKeyJSON []byte
	var cnt int
	for rows.Next() {
		cnt++
		err = rows.Scan(&username, &pwd,&encPrivateKeyJSON)
        checkErr(err)
	}
	rows.Close()
	db.Close()
	
	if cnt == 0 {
		return nil
	}

	privateKeyJSON := AESdecryptdata([]byte(password),encPrivateKeyJSON)
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
	returnUser.Name = username
	returnUser.Password = password
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

// PUBLISHER

// SQLiteInitPshDB ´sdd
func SQLiteInitPshDB(dbpath string) {
	 db, err := sql.Open("sqlite3", dbpath)
		checkErr(err)		
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
	checkErr(err)
	guid := SHA1HashString(jsonBlob)
	// insert
    stmt, err := db.Prepare("INSERT INTO publish(data, created, identifier ) values(?,?,?)")
    checkErr(err)
	res, err := stmt.Exec(jsonBlob, time.Now(),SHA1HashString(jsonBlob))
    checkErr(err)
    id, err := res.LastInsertId()
	checkErr(err)
	stmt.Close()
	log.Println(id)
	return guid
}


func SQLiteGetALLPublishData(dbpath string) {
	
	db, err := sql.Open("sqlite3", dbpath)
    checkErr(err)
	// todo: has to be better
	rows, err := db.Query("Select uid,data,identifier from publish")
	checkErr(err)
	var uid int
	var data []byte
	var identifier string
	var cnt int
	for rows.Next() {
		cnt++
		err = rows.Scan(&uid, &data,&identifier)
		checkErr(err)
		log.Printf("%d %s",uid,identifier)
	}
	log.Println(cnt)
	rows.Close()
	db.Close()
}


func SQLiteGetPublishedData(dbpath string, hash string) []byte {
	
	db, err := sql.Open("sqlite3", dbpath)
    checkErr(err)
    // todo: has to be better
    SQLQuery := "Select data from publish where identifier = '"+hash+"'"
    log.Debug(SQLQuery)
	rows, err := db.Query(SQLQuery)
	checkErr(err)
	var data []byte
	for rows.Next() {
        err = rows.Scan(&data)
		checkErr(err)
	}
	rows.Close()
    db.Close()
    return data
}
