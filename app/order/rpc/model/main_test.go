package model

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DbDriver = "postgres"
	DbSource = "postgresql://postgres:postgres@localhost:35434/shop?sslmode=disable"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(DbDriver, DbSource)
	if err != nil {
		log.Fatalln("cannot connect to db :", err)
	}
	testQueries = New(testDB)
	log.Println("connect db success....")
	// m.Run() 返回一个退出的代码，告诉我们测试是否通过
	// 使用 os.Exit() 将测试的结果报告给测试运行程序
	os.Exit(m.Run())
}
