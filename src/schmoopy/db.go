package schmoopy

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	"code.google.com/p/go-sqlite/go1/sqlite3"
)

func InitializeDb(
	dbFilename string,
) error {
	if _, err := os.Stat(dbFilename); err == nil {
		log.Fatal(fmt.Errorf("File already exists: %v", dbFilename))
	}

	c, err := sqlite3.Open(dbFilename)
	if err != nil {
		return err
	}

	sql := "CREATE TABLE schmoopys(" +
		"schmoopy STRING," +
		"imageUrl STRING," +
		"PRIMARY KEY (schmoopy, imageUrl))"

	if err = c.Exec(sql); err != nil {
		return err
	}

	if err = c.Close(); err != nil {
		return err
	}

	return nil
}

type schmoopy struct {
	Name      string
	ImageUrls map[string]struct{}
}

func (s *schmoopy) RandomImageUrl() string {
	urls := []string{}
	for url, _ := range s.ImageUrls {
		urls = append(urls, url)
	}
	return urls[rand.Intn(len(urls))]
}

func dbAddSchmoopy(conn *sqlite3.Conn, name string, imageUrl string) error {
	args := sqlite3.NamedArgs{"$schmoopy": name, "$imageUrl": imageUrl}
	sql := "INSERT OR IGNORE INTO schmoopys(schmoopy, imageUrl) " +
		"VALUES ($schmoopy, $imageUrl)"
	return conn.Exec(sql, args)
}

func dbRemoveSchmoopy(conn *sqlite3.Conn, name string, imageUrl string) error {
	args := sqlite3.NamedArgs{"$schmoopy": name, "$imageUrl": imageUrl}
	sql := "DELETE FROM schmoopys WHERE schmoopy = $schmoopy AND imageUrl = $imageUrl"
	return conn.Exec(sql, args)

}

func dbFetchSchmoopys(conn *sqlite3.Conn, names []string) (map[string]*schmoopy, error) {
	sql := "SELECT schmoopy, imageUrl FROM schmoopys"
	args := make([]interface{}, len(names))

	if len(names) > 0 {
		params := make([]string, len(names))
		for idx, name := range names {
			params[idx] = "?"
			args[idx] = name
		}
		sql += " WHERE schmoopy IN (" + strings.Join(params, ", ") + ")"
	}
	schmoopys := map[string]*schmoopy{}

	rows, err := conn.Query(sql, args...)

	if err == io.EOF {
		return schmoopys, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	var name string
	var imageUrl string

	for {
		err = rows.Scan(&name, &imageUrl)

		if err != nil {
			return nil, err
		}

		sch, ok := schmoopys[name]
		if !ok {
			sch = &schmoopy{
				Name:      name,
				ImageUrls: map[string]struct{}{},
			}
			schmoopys[name] = sch
		}
		sch.ImageUrls[imageUrl] = struct{}{}

		err = rows.Next()

		// EOF signifies end of the query
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return schmoopys, nil
}
