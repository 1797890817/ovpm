package ovpm_test

import (
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/cad/ovpm"
)

func TestCreateNewUser(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")
	server, _ := ovpm.GetServerInstance()

	// Prepare:
	username := "testUser"
	password := "testPasswd1234"

	// Test:
	user, err := ovpm.CreateNewUser(username, password)
	if err != nil {
		t.Fatalf("user can not be created: %v", err)
	}

	// Is user nil?
	if user == nil {
		t.Fatalf("user is expected to be 'NOT nil' but it is 'nil' %+v", user)
	}

	// Is user empty?
	if *user == (ovpm.DBUser{}) {
		t.Fatalf("user is expected to be 'NOT EMPTY' but it is 'EMPTY' %+v", user)
	}

	// Is user acutally exist in the system?
	user2, err := ovpm.GetUser(username)
	if err != nil {
		t.Fatalf("user can not be retrieved: %v", err)
	}

	// Are users the same?
	if !areUsersEqual(user, user2) {
		t.Fatalf("users are expected to be 'NOT TO DIFFER' but they are 'DIFFERENT'")
	}

	// Is user's server serial number correct?
	if !server.CheckSerial(user.ServerSerialNumber) {
		t.Fatalf("user's ServerSerialNumber is expected to be 'CORRECT' but it is 'INCORRECT' instead %+v", user)
	}

	// Does User interface work properly?
	if user.GetUsername() != user.Username {
		t.Errorf("user.GetUsername() is expected to return '%s' but it returns '%s' %+v", user.Username, user.GetUsername(), user)
	}
	if user.GetServerSerialNumber() != user.ServerSerialNumber {
		t.Errorf("user.GetServerSerialNumber() is expected to return '%s' but it returns '%s' %+v", user.ServerSerialNumber, user.GetServerSerialNumber(), user)
	}
	if user.GetCert() != user.Cert {
		t.Errorf("user.GetCert() is expected to return '%s' but it returns '%s' %+v", user.Cert, user.GetCert(), user)
	}

}

func TestUserPasswordCorrect(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")

	// Prepare:
	initialPassword := "g00dp@ssW0rd9"
	user, _ := ovpm.CreateNewUser("testUser", initialPassword)

	// Test:

	// Is user created with the correct password?
	if !user.CheckPassword(initialPassword) {
		t.Fatalf("user's password must be '%s', but CheckPassword fails +%v", initialPassword, user)
	}
}

func TestUserPasswordReset(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")

	// Prepare:
	initialPassword := "g00dp@ssW0rd9"
	user, _ := ovpm.CreateNewUser("testUser", initialPassword)

	// Test:

	// Reset user's password.
	newPassword := "@n0th3rP@ssw0rd"
	user.ResetPassword(newPassword)

	// Is newPassword set correctly?
	if !user.CheckPassword(newPassword) {
		t.Fatalf("user's password must be '%s', but CheckPassword fails +%v", newPassword, user)
	}

	// Is initialPassword is invalid now?
	if user.CheckPassword(initialPassword) {
		t.Fatalf("user's password must be '%s', but CheckPassword returns true for the old password '%s' %+v", newPassword, initialPassword, user)
	}
}

func TestUserDelete(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")

	// Prepare:
	username := "testUser"
	user, _ := ovpm.CreateNewUser(username, "1234")

	// Test:

	// Is user fetchable?
	_, err := ovpm.GetUser(username)
	if err != nil {
		t.Fatalf("user '%s' expected to be 'EXIST', but we failed to get it %+v: %v", username, user, err)
	}

	// Delete the user.
	err = user.Delete()

	// Is user deleted?
	if err != nil {
		t.Fatalf("user is expected to be deleted, but we got this error instead: %v", err)
	}

	// Is user now fetchable? (It shouldn't be.)
	u1, err := ovpm.GetUser(username)
	if err == nil {
		t.Fatalf("user '%s' expected to be 'NOT EXIST', but it does exist instead: %v", username, user)
	}

	// Is user nil?
	if u1 != nil {
		t.Fatalf("not found user should be 'nil' but it's '%v' instead", u1)
	}
}

func TestUserGet(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")

	// Prepare:
	username := "testUser"
	user, _ := ovpm.CreateNewUser(username, "1234")

	// Test:
	// Is user fetchable?
	fetchedUser, err := ovpm.GetUser(username)
	if err != nil {
		t.Fatalf("user should be fetchable but instead we got this error: %v", err)
	}

	// Is fetched user same with the created one?
	if !areUsersEqual(user, fetchedUser) {
		t.Fatalf("fetched user should be same with the created user but it differs")
	}

}

func TestUserGetAll(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")
	count := 5

	// Prepare:
	var users []*ovpm.DBUser
	for i := 0; i < count; i++ {
		username := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("password%d", i)
		user, _ := ovpm.CreateNewUser(username, password)
		users = append(users, user)
	}

	// Test:
	// Get all users.
	fetchedUsers, err := ovpm.GetAllUsers()

	// Is users are fetchable.
	if err != nil {
		t.Fatalf("users are expected to be fetchable but we got this error instead: %v", err)
	}

	// Is the fetched user count is correct?
	if len(fetchedUsers) != count {
		t.Fatalf("fetched user count is expected to be '%d', but it is '%d' instead", count, len(fetchedUsers))
	}

	// Are returned users same with the created ones?
	for i, user := range fetchedUsers {
		if !areUsersEqual(user, users[i]) {
			t.Fatalf("user %s[%d] is expected to be 'SAME' with the created one of it, but they are 'DIFFERENT'", user.GetUsername(), i)
		}
	}
}

func TestUserRenew(t *testing.T) {
	// Initialize:
	ovpm.Testing = true
	ovpm.SetupDB("sqlite3", ":memory:")
	defer ovpm.CeaseDB()
	ovpm.Init("localhost", "")

	// Prepare:
	user, _ := ovpm.CreateNewUser("user", "1234")

	// Test:
	// Re initialize the server.
	ovpm.Init("example.com", "3333") // This causes implicit Renew() on every user in the system.

	// Fetch user back.
	fetchedUser, _ := ovpm.GetUser(user.GetUsername())

	// Aren't users certificates different?
	if user.Cert == fetchedUser.Cert {
		t.Fatalf("fetched user's certificate is expected to be 'DIFFERENT' from the created user (since server is re initialized), but it's 'SAME' instead")
	}
}

// areUsersEqual compares given users and returns true if they are the same.
func areUsersEqual(user1, user2 *ovpm.DBUser) bool {
	if user1.Cert != user2.Cert {
		logrus.Info("Cert %v != %v", user1.Cert, user2.Cert)
		return false
	}
	if user1.Username != user2.Username {
		logrus.Infof("Username %v != %v", user1.Username, user2.Username)
		return false
	}
	if user1.Password != user2.Password {
		logrus.Infof("Password %v != %v", user1.Password, user2.Password)
		return false
	}
	if user1.ServerSerialNumber != user2.ServerSerialNumber {
		logrus.Infof("ServerSerialNumber %v != %v", user1.ServerSerialNumber, user2.ServerSerialNumber)
		return false
	}
	logrus.Infof("users are the same!")
	return true

}
