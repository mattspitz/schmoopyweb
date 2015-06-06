package schmoopy

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"fmt"
	"io"
	"log"
	"os"
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
	name      string
	imageUrls map[string]struct{}
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
	args := sqlite3.NamedArgs{}
	if len(names) > 0 {
		// TODO join on ,
		sql += " WHERE schmoopy = $name"
		args["$name"] = names[0]
	}
	schmoopys := map[string]*schmoopy{}

	rows, err := conn.Query(sql, args)

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
				name:      name,
				imageUrls: map[string]struct{}{},
			}
			schmoopys[name] = sch
		}
		sch.imageUrls[imageUrl] = struct{}{}

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
