package sql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//ClientConn ...
type ClientConn struct {
	address  string
	usrname  string
	password string
	schema   string
	db       *sql.DB
}

//NewClientConn ...
func NewClientConn(address, usrname, password, schema string) (*ClientConn, error) {
	cc := &ClientConn{
		address:  address,
		usrname:  usrname,
		password: password,
		schema:   schema,
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", usrname, password, address, schema) //"user:password@/dbname"

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	cc.db = db

	return cc, nil //db.Ping()
}

//Close ...
func (cc *ClientConn) Close() {
	cc.db.Close()
}

//QueryStoreURL ...
func (cc *ClientConn) QueryStoreURL(resourceID string) (string, string, error) {
	rows, err := cc.db.Query("SELECT base_path, file FROM media_store WHERE resource_id = ?", resourceID)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	var basePath, file string
	for rows.Next() {
		if err := rows.Scan(&basePath, &file); err != nil {
			return "", "", err
		}
		return basePath, file, nil
	}

	return "", "", errors.New("no result")
}

//QuerySteamingURL ...
func (cc *ClientConn) QuerySteamingURL(resourceID string) (string, string, error) {
	rows, err := cc.db.Query("SELECT base_path, endpoint FROM media_stream WHERE resource_id = ?", resourceID)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	var basePath, endpoint string
	for rows.Next() {
		if err := rows.Scan(&basePath, &endpoint); err != nil {
			return "", "", err
		}
		return basePath, endpoint, nil
	}

	return "", "", errors.New("no result")
}

//AddStreamingURL ...
func (cc *ClientConn) AddStreamingURL(resourceID, basePath, endpoint string) error {
	stmt, err := cc.db.Prepare("INSERT INTO media_stream (`resource_id`, `base_path`, `endpoint`) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(resourceID, basePath, endpoint)
	return err
}

//UpdateStreamingURL ...
func (cc *ClientConn) UpdateStreamingURL(resourceID, basePath, endpoint string) error {
	stmt, err := cc.db.Prepare("UPDATE media_stream SET base_path=?, endpoint=? WHERE resource_id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(basePath, endpoint, resourceID)
	return err
}

//DeleteStreamingURL ...
func (cc *ClientConn) DeleteStreamingURL(resourceID string) error {
	stmt, err := cc.db.Prepare("DELETE FROM media_stream WHERE resource_id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(resourceID)
	return err
}
