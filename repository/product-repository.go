package repository

import (
	"database/sql"
	"fmt"
	"time"
	"viewscounter/entity"
)

type ProductRepository interface {
	InsertResViews(b []entity.ResInfo)
	InsertMediaViews(b []entity.MediaInfo)
	InsertRpAccViews(b []entity.RpAccInfo)
	AllProduct() []entity.Product
}

type productConnection struct {
	connection *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productConnection{
		connection: db,
	}
}

//Increment of views
func (db *productConnection) InsertResViews(b []entity.ResInfo) {
	for _, val := range b {
		_, err := db.connection.Exec(`UPDATE tbl_dk_resource SET "ResViewCnt" = "ResViewCnt" + 1 WHERE "ResGuid"=$1 AND "ResRegNo"=$2;`, val.ProductGuid, val.ResNo)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("Updated")
}

func (db *productConnection) InsertMediaViews(b []entity.MediaInfo) {
	for _, val := range b {
		_, err := db.connection.Exec(`UPDATE tbl_me_media SET "MediaViewCnt" = "MediaViewCnt" + 1 WHERE "MediaGuid"=$1;`, val.MediaGuid)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("Updated")
}

func (db *productConnection) InsertRpAccViews(b []entity.RpAccInfo) {
	for _, val := range b {
		_, err := db.connection.Exec(`UPDATE tbl_dk_rp_acc SET "RpAccViewCnt" = "RpAccViewCnt" + 1 WHERE "RpAccGuid"=$1 AND "RpAccRegNo"=$2;`, val.RpAccGuid, val.RpAccNo)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("Updated")
}

func (db *productConnection) AllProduct() []entity.Product {
	var (
		product  entity.Product
		products []entity.Product
	)
	rows, err := db.connection.Query(`SELECT dk."ResGuid", dk."ResName", dk."ResId" FROM tbl_dk_resource as dk, ez_mod_time e WHERE dk."ModifiedDate" > e."ModifiedDate";`)
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&product.ResGuid, &product.ResName, &product.ResId)
		products = append(products, product)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	fmt.Println(time.Now().Format("2006-01-02"))
	_, err = db.connection.Exec(`UPDATE ez_mod_time SET "ModifiedDate" = $1;`, time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Print(err.Error())
	}
	defer rows.Close()
	fmt.Println("Synced")
	return products
}
