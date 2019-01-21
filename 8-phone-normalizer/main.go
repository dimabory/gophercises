package main

import (
	"fmt"
	phonedb "github.com/dimabory/gophercises/8-phone-normalizer/db"
	_ "github.com/lib/pq"
	"regexp"
)

// see configuration in docker-compose.yml
const (
	host          = "localhost"
	port          = 3254
	user          = "Dmytr0"
	password      = "123123"
	dbname        = "gophercises"
	defaultDriver = "postgres"
)

var connectionString string

func init() {
	connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	must(phonedb.Migrate(defaultDriver, connectionString))
}

func main() {
	db, err := phonedb.Open(defaultDriver, connectionString)
	must(err)

	defer db.Close()

	must(db.Seed())

	phones, err := db.AllPhones()
	must(err)

	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println(" --> ", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.ID))
			} else {
				p.Number = number
				must(db.UpdatePhone(&p))
			}
			continue
		}

		fmt.Println("No changes required")
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")

	return re.ReplaceAllString(phone, "")
}

//func normalize(phone string) string {
//	var buf bytes.Buffer
//
//	for _, ch := range phone {
//		if ch >= '0' && ch <= '9' {
//			buf.WriteRune(ch)
//		}
//	}
//
//	return buf.String()
//}
