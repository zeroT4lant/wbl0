package database

import "fmt"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// "postgresql://postgres:Amogus228!@localhost:5432/wbdb"
func (c Config) connStr() string {
	return fmt.Sprintf("host=%s port=%s username=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
}
