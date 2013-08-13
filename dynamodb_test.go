package dynamodb_test

import (
	"log"
	"testing"
	"time"

	_ "github.com/eikeon/dynamodb/driver/dynamo"
	_ "github.com/eikeon/dynamodb/driver/memory"

	"github.com/eikeon/dynamodb"
)

type Fetch struct {
	URL string `db:"HASH"`
}

var FETCH *Fetch

func TestCreateTablePutItemScanDeleteTable(t *testing.T) {
	for _, driverName := range []string{"memory", "dynamo"} {
		d, err := dynamodb.Open(driverName)
		if err != nil {
			t.Fatal(err)
		}

		d.Register("fetch", (*Fetch)(nil))

		if err := d.CreateTable("fetch"); err != nil {
			t.Error(err)
		}

		for {
			if description, err := d.DescribeTable("fetch"); err != nil {
				t.Error(err)
			} else {
				log.Println(description.Table.TableStatus)
				if description.Table.TableStatus == "ACTIVE" {
					break
				}
			}
			time.Sleep(time.Second)
		}

		if err := d.PutItem("fetch", &Fetch{"http://localhost/"}); err != nil {
			t.Error(err)
		}

		if response, err := d.Scan("fetch"); err != nil {
			t.Error(err)
		} else {
			items := response.GetItems()
			for i := 0; i < response.GetCount(); i++ {
				item := items[i]
				log.Println("item:", item)
			}
		}

		time.Sleep(10 * time.Second)

		if err := d.DeleteTable("fetch"); err != nil {
			t.Error(err)
		}
	}
}
