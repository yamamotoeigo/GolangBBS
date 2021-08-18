package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Comment struct {
	gorm.Model
	Name string
	Body string
}

/**
DB初期化
*/
func dbInit() {
	db, err := gorm.Open("sqlite3", "./comment.sqlite3")
	if err != nil {
		panic("DB開けず！(dbInit)")
	}
	db.AutoMigrate(&Comment{})
	defer db.Close()
}

/**
DB追加
*/
func dbInsert(name string, body string) {
	db, err := gorm.Open("sqlite3", "./comment.sqlite3")
	if err != nil {
		panic("DB開けず(dbInsert)")
	}
	db.Create(&Comment{Name: name, Body: body})
	defer db.Close()
}

/**
DB更新
*/
func dbUpdate(id int, name string, body string) {
	db, err := gorm.Open("sqlite3", "./Comment.sqlite3")
	if err != nil {
		panic("DB開けず(dbUpdate)")
	}
	var comment Comment
	db.First(&comment, id)
	comment.Name = name
	comment.Body = body
	db.Save(&comment)
	db.Close()
}

/**
DB削除
*/
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "./comment.sqlite3")
	if err != nil {
		panic("DB開けず(dbDelete")
	}
	var comment Comment
	db.First(&comment, id)
	db.Delete(&comment)
	db.Close()
}

/**
DB全取得
*/
func dbGetAll() []Comment {
	db, err := gorm.Open("sqlite3", "./comment.sqlite3")
	if err != nil {
		panic("DB開けず(dbGetAll)")
	}
	var comments []Comment
	db.Order("created_at desc").Find(&comments)
	db.Close()
	return comments
}

/**
DB一つ取得
*/
func dbGetOne(id int) Comment {
	db, err := gorm.Open("sqlite3", "./comment.sqlite3")
	if err != nil {
		panic("DB開けず(dbGetOne)")
	}
	var comment Comment
	db.First(&comment, id)
	db.Close()
	return comment
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		comments := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"comments": comments,
		})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		body := ctx.PostForm("body")
		dbInsert(name, body)
		ctx.Redirect(302, "/")
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		comment := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"comment": comment})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		name := ctx.PostForm("name")
		body := ctx.PostForm("body")
		dbUpdate(id, name, body)
		ctx.Redirect(302, "/")
	})

	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		comment := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"comment": comment})
	})

	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()
}
