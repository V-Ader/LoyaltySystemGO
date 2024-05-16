package client

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUsers(dbConnection *sql.DB) ([]Client, error) {
	results, err := dbConnection.Query("SELECT * FROM clients")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	users := []Client{}
	for results.Next() {
		var user Client
		err = results.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetClientById(dbConnection *sql.DB, context *gin.Context) (*Client, error) {
	id := context.Param("id")
	query := "SELECT * FROM clients WHERE id = $1"
	results, err := dbConnection.Query(query, id)
	if err != nil {
		fmt.Println("asd")
		return nil, err
	}
	defer results.Close()

	user := &Client{}
	if results.Next() {
		err = results.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("record not found")
	}
	return user, nil
}

func ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
	query := "INSERT INTO clients (id, name, email) VALUES (nextval('client_seq'), $1, $2)"
	_, err := dbConnection.Exec(query, context.Query("name"), context.Query("email"))
	return err
}

func ExecuteDelte(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM clients where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}

// func ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
// 	query := "INSERT INTO users (id, name) VALUES (nextval('user_seq'), $1)"
// 	_, err := dbConnection.Exec(query, context.Query("name"))
// 	return err
// }
