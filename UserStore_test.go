package main

import "testing"

func Test_UserStoreConstructor(t *testing.T) {

	dbpath := "./hier.db"
	tUserStore := NewUserstore(dbpath)

	if tUserStore.userCount != 0 {
		t.Errorf("Data missmatch %v != %d", tUserStore.userCount, 0)
	}
}



func Test_UserStoreAddition(t *testing.T) {
	dbpath := "./hier.db"
	tUserStore := NewUserstore(dbpath)
	tUserStore.newUser("testuser", "password")

	if tUserStore.userCount != 1 {
		t.Errorf("missmatch: tUserStore.userCount %v != %d", tUserStore.userCount, 1)
	}
}

func Test_UserStoreGetKeyPair(t *testing.T) {

	dbpath := "./hier.db"
	tUserStore := NewUserstore(dbpath)
	tUserStore.newUser("testuser", "password")
	tprivKey := tUserStore.getKeyPair("testuser")
	if tprivKey == nil {
		t.Errorf("missmatch: tprivKey %v == nil", tprivKey )
	}


}