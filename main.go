package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

// DBマイグレート
func db_init() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("We cant open DB (db_init)")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

// DBに追加
func db_create(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("we cant open DB (db_create)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

// DBの要素を全取得
func db_get_all() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("we cant open DB (db_get_all)")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()

	return todos
}

// 指定されたidの要素をDBから取得
func db_get_one(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("we cant open DB (db_get_one)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()

	return todo
}

// DBを更新
func db_update(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("we cant open DB (db_update)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

// DBから要素を削除
func db_delete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("we cant open DB (db_delete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	db_init()

	// Index
	router.GET("/", func(ctx *gin.Context) {
		todos := db_get_all()
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	// Create
	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		db_create(text, status)
		ctx.Redirect(302, "/")
	})

	// Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := db_get_one(id)
		ctx.HTML(200, "index.html", gin.H{
			"todo": todo,
		})
	})

	// Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		db_update(id, text, status)
		ctx.Redirect(302, "/")
	})

	// 削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := db_get_one(id)
		ctx.HTML(200, "delete.html", gin.H{
			"todo": todo,
		})
	})

	// Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		db_delete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()
}
