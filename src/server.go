package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/labstack/echo"
    mw "github.com/labstack/echo/middleware"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    "data"
)

var con *sql.DB

func createPosts(con *sql.DB) {
    var post1, post2, post3 data.Post

    post1.Title = "Post 1"
    post1.Body = "<p>Lorem ipsum, dolor sit amet</p>"
    post1.Author = "Reset Reboot"

    post1.Save(con)

    post2.Title = "Echo, Golang framework"
    post2.Body = "<p>Parapapapapa, the rapper, wawawawa</p>"
    post2.Author = "Reset Reboot"

    post2.Save(con)

    post3.Title = "Python still rocks, anyway"
    post3.Body = "<p>import antigravity<br /># Fly!</p>"
    post3.Author = "Reset Reboot"

    post3.Save(con)
}

// Handler
func hello(c *echo.Context) error {
    var text, temp string

    posts := data.GetPostList(con, 0, 10)

    for _, post := range posts {
        temp = fmt.Sprintf("%v", post)
        text = text + "<br />" + temp
    }

    return c.HTML(http.StatusOK, "<html><body><h1>Blog posts:</h1>" + text + "</body></html>")
}

func main() {
    var err error
    // Echo instance
    e := echo.New()
    e.SetDebug(true)

    // Middleware
    e.Use(mw.Logger())
    e.Use(mw.Recover())

    con, err = sql.Open("mysql", "golang:golang@tcp(localhost:3306)/golang")
    if err != nil {
        log.Fatal(err)
    }

    data.InitDatabase(con)

    // createPosts(con)

    // Routes
    e.Get("/", hello)

    // Start server
    e.Run(":1323")
}
