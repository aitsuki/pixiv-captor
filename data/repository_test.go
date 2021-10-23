package data

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var testData = Illust{
	ID:  "1",
	R18: false,
}

var testDataList = []Illust{
	{
		ID:          "1",
		R18:         false,
		Description: "碧蓝航线, 赤诚",
	},
	{
		ID:          "2",
		R18:         true,
		Description: "加贺, 碧蓝航线",
	},
	{
		ID:          "3",
		R18:         true,
		Description: "爱宕, 碧蓝航线",
	},
}

func openTestDB(t *testing.T) *sql.DB {
	temp, err := os.CreateTemp("", "test.db")
	if err != nil {
		t.Fatal("Failed to create test temp file:", err)
	}
	temp.Close()
	db, err := sql.Open("sqlite3", temp.Name())
	// db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	return db
}

func Test_Prepare(t *testing.T) {
	repo := NewIllustRepository(openTestDB(t))
	err := repo.Prepare()
	if err != nil {
		t.Fatal("Failed to prepare repository:", err)
	}
}

func Test_IsExists(t *testing.T) {
	repo := NewIllustRepository(openTestDB(t))
	repo.Prepare()
	exists := repo.IsExists(testData.ID)
	if exists {
		t.Fatalf("IsExists() = %v, want %v", exists, false)
	}

	repo.Save(&testData)

	exists = repo.IsExists(testData.ID)
	if !exists {
		t.Fatalf("IsExists() = %v, want %v", exists, true)
	}
}

func Test_Save(t *testing.T) {
	repo := NewIllustRepository(openTestDB(t))
	repo.Prepare()
	err := repo.Save(&testData)
	if err != nil {
		t.Fatal("Failed to save data:", err)
	}
	exists := repo.IsExists(testData.ID)
	if !exists {
		t.Fatalf("IsExists() = %v, want %v", exists, true)
	}
}

func Test_GetByID(t *testing.T) {
	repo := NewIllustRepository(openTestDB(t))
	repo.Prepare()
	repo.Save(&testData)
	illust, err := repo.GetByID(testData.ID)
	if err != nil {
		t.Fatal("Failed to get data by ID:", err)
	}
	if illust.ID != testData.ID {
		t.Fatalf("illust.ID = %v, want %v", testData.ID, illust.ID)
	}
}

func Test_GetRandom(t *testing.T) {
	repo := NewIllustRepository(openTestDB(t))
	repo.Prepare()
	for _, testData := range testDataList {
		repo.Save(&testData)
	}

	type args struct {
		r18   int
		limit int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"one rows r18", args{1, 1}, 1},
		{"tow rows r18", args{1, 2}, 2},
		{"one rows not r18", args{0, 1}, 1},
		{"no result", args{-1, 1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetRandom(tt.args.r18, tt.args.limit)
			if err != nil {
				t.Error("Failed to get random data:", err)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GetRandom(%v, %v).len(rows) = %v, want %v", tt.args.r18, tt.args.limit, len(got), tt.want)
			}
		})
	}
}

func Test_Search(t *testing.T) {
	repo := NewIllustRepository(openTestDB(t))
	repo.Prepare()
	for _, testData := range testDataList {
		repo.Save(&testData)
	}

	type args struct {
		r18   int
		q     string
		limit int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"one rows r18", args{1, "加贺", 1}, 1},
		{"tow rows r18", args{1, "碧蓝", 2}, 2},
		{"one rows not r18", args{0, "赤诚", 1}, 1},
		{"no result", args{1, "赤诚", 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Search(tt.args.r18, tt.args.q, tt.args.limit)
			if err != nil {
				t.Error("Failed to search data:", err)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Search(%v, %v, %v).len(rows) = %v, want %v", tt.args.r18, tt.args.q, tt.args.limit, len(got), tt.want)
			}
		})
	}
}
