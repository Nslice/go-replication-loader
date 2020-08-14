package mssql

import (
	"fmt"
	"strings"

	"github.com/sergeyzalunin/go-replication-loader/argsp"
)

// NewConnectionString is a method factory to construct connection string
func NewConnectionString(args argsp.ArgumentOptions) string {
	conn := &connectionStringBuilder{[]string{}}
	result := conn.server(args.DbDataSource).
		database(args.DatabaseName).
		trustedConnection(args.TrustedConnection).
		connectionTimeout(args.ConnectionTimeout).
		userID(args.DatabaseUserID).
		password(args.DatabasePassword).
		build()
	return result
}

type connectionStringBuilder struct {
	connection []string
}

func (cs *connectionStringBuilder) server(server string) *connectionStringBuilder {
	if strings.TrimSpace(server) != "" {
		serverConn := fmt.Sprintf("Server=%s", server)
		cs.connection = append(cs.connection, serverConn)
	}
	return cs
}

func (cs *connectionStringBuilder) database(db string) *connectionStringBuilder {
	if strings.TrimSpace(db) != "" {
		dbConn := fmt.Sprintf("Database=%s", db)
		cs.connection = append(cs.connection, dbConn)
	}
	return cs
}

func (cs *connectionStringBuilder) trustedConnection(trust bool) *connectionStringBuilder {
	dbConn := fmt.Sprintf("Trusted_Connection=%t", trust)
	cs.connection = append(cs.connection, dbConn)
	return cs
}

func (cs *connectionStringBuilder) connectionTimeout(timeout int) *connectionStringBuilder {
	dbConn := fmt.Sprintf("Connection Timeout=%d", timeout)
	cs.connection = append(cs.connection, dbConn)
	return cs
}

func (cs *connectionStringBuilder) userID(user string) *connectionStringBuilder {
	if strings.TrimSpace(user) != "" {
		userConn := fmt.Sprintf("User Id=%s", user)
		cs.connection = append(cs.connection, userConn)
	}
	return cs
}

func (cs *connectionStringBuilder) password(password string) *connectionStringBuilder {
	if strings.TrimSpace(password) != "" {
		passConn := fmt.Sprintf("Password=%s", password)
		cs.connection = append(cs.connection, passConn)
	}
	return cs
}

func (cs *connectionStringBuilder) build() string {
	result := strings.Join(cs.connection, "; ")
	return result
}
