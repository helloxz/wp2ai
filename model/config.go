package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/glebarez/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var WP *gorm.DB

func InitWPDB() {
	host := viper.GetString("wordpress.db_host")
	name := viper.GetString("wordpress.db_name")
	username := viper.GetString("wordpress.db_username")
	password := viper.GetString("wordpress.db_password")
	dsn := username + ":" + password + "@tcp(" + host + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
	// dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	if WP == nil {
		WP, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			// log.Fatalf("failed to connect database: %v", err)
			// 写入日志
			// 打印错误消息
			fmt.Printf("failed to connect database: %v\n", err)
		} else {
			// 打印连接成功
			fmt.Printf("WordPress database connection succeeded!\n")
		}
	}

}

// 全局数据库连接
var DB *gorm.DB

// 初始化数据库连接
func InitDB() {
	var err error
	//如果数据库文件不存在，会自动创建
	DB, err = gorm.Open(sqlite.Open("data/db/wp2db.db3"), &gorm.Config{})
	//如果出现错误，抛出错误并终止执行
	if err != nil {
		// 打印错误
		fmt.Println("Database connection failed!")
		// 写入日志
		// helper.WriteLog(fmt.Sprintf("%s", err))
		// 终止执行
		os.Exit(1)
	} else {
		// 显式启用 WAL 模式
		if err := DB.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
			log.Fatalf("Failed to enable WAL mode: %v", err)
		}
		// 调整 SQLite 参数
		DB.Exec("PRAGMA synchronous=NORMAL;")
		//自动迁移数据库表结构
		// 迁移多个表
		err = DB.AutoMigrate(&Post{})
		if err != nil {
			fmt.Println("failed to migrate database!")
			// 写入错误日志
			// helper.WriteLog(fmt.Sprintf("%s", err))
		}
		fmt.Print("Database connection succeeded!\n")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("failed to get database connection pool!")
		os.Exit(1)
	}
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(5)
	// 设置连接的最大生命周期
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

}

var VecDB *sql.DB

func InitVecDB() {
	sqlite_vec.Auto()
	var err error
	// Open persistent database
	VecDB, err = sql.Open("sqlite3", "data/db/embedding.db3")
	if err != nil {
		log.Fatal(err)
	}
	// defer VecDB.Close()

	var sqliteVersion, vecVersion string
	err = VecDB.QueryRow("select sqlite_version(), vec_version()").Scan(&sqliteVersion, &vecVersion)
	if err != nil {
		log.Fatal("failed to prepare statement!")
	}

	fmt.Printf("sqlite_version=%s, vec_version=%s\n", sqliteVersion, vecVersion)

	// Create table "items" with specified fields
	_, err = VecDB.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY,
			post_id TEXT NOT NULL UNIQUE,
			embedding BLOB NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("failed to create table items!")
	}
}
