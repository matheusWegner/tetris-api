package main

/*import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// Configurando a conexão com o MySQL
	dsn := "root:root@tcp(localhost:3306)/tetris"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Criando a tabela `users`
	createTableUsersSql := `
        CREATE TABLE users (
            id INT NOT NULL AUTO_INCREMENT,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
            PRIMARY KEY (id)
        );
    `
	_, err = db.Exec(createTableUsersSql)
	if err != nil {
		log.Fatal(err)
	}
}

func getUsers(db *sql.DB) ([]*User, error) {
	// Executando a consulta SQL
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterando sobre os resultados
	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Retornando os resultados
	return users, nil
}

func getUser(db *sql.DB, id int) (*User, error) {
	// Executando a consulta SQL
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	// Deserializando o resultado
	user := &User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	// Retornando o resultado
	return user, nil
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}
*/
