package data

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var illustTestData = Illust{
	ID:          "80971280",
	Title:       "可可萝的高抬腿练习",
	Description: `看了p站画师まんなく（来自妖梦喵的提示）的画，有所感触，也作练习一张。每月发布当月的无码图，感谢支持！`,
	AuthorID:    "9751291",
	Author:      "大猫板蓝根",
	Account:     "ex_azusa",
	R18:         true,
	CreateDate:  time.Date(2020, 4, 21, 23, 29, 27, 0, time.UTC),
	UploadDate:  time.Date(2020, 4, 21, 23, 29, 27, 0, time.UTC),
	Tags:        []string{"可可萝", "Kokkoro", "开腿", "公主连结Re:Dive", "コッコロ", "魅惑のふともも", "魅惑的大腿"},
	Pages: []IllustPage{
		{
			IllustID: "80971280",
			P:        0,
			Width:    1200,
			Height:   1405,
			Thumb:    `https://i.pximg.net/c/128x128/img-master/img/2020/04/22/08/29/27/80971280_p0_square1200.jpg`,
			Small:    `https://i.pximg.net/c/540x540_70/img-master/img/2020/04/22/08/29/27/80971280_p0_master1200.jpg`,
			Regular:  `https://i.pximg.net/c/540x540_70/img-master/img/2020/04/22/08/29/27/80971280_p0_master1200.jpg`,
			Original: `https://i.pximg.net/img-original/img/2020/04/22/08/29/27/80971280_p0.jpg`,
		},
	},
}

func TempFilename(t *testing.T) string {
	f, err := ioutil.TempFile("", "go-sqlite3-test-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func Test_save_and_find(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)
	db, err := sql.Open("sqlite3", tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	dao, err := NewIllustDao(db)
	if err != nil {
		t.Fatal("Failed to create_dao", err)
	}

	err = dao.Save(&illustTestData)
	if err != nil {
		t.Fatal(err)
	}

	row := dao.db.QueryRow("select tag from illust_tags where illust_id = illusts.id")
	s := ""
	err = row.Scan(&s)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(s)
}
