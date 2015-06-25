package data

import (
    "fmt"
    "log"
    "time"
    "database/sql"
)

const human_time_layout = "Mon Jan 2 15:04:05"
const time_layout = "2006-01-02 15:04:05"

const database_command = "CREATE TABLE IF NOT EXISTS post (id INTEGER AUTO_INCREMENT PRIMARY KEY, title VARCHAR (255), body TEXT, author VARCHAR(120), date DATETIME)"

func InitDatabase(con *sql.DB) {
    _, err := con.Exec(database_command)

    if err != nil {
        log.Print(err)
        return
    }
}

// Post struct and methods

type Post struct {
    ID       int
    Title    string
    Body     string
    Author   string
    PostDate time.Time
    isSaved  bool
}

func (p *Post) Save(con *sql.DB) {
    var statement string
    var err error
    var res sql.Result

    if !p.isSaved {
        statement = "INSERT INTO post(title, body, author, date) VALUES (?,?,?,?)"

    } else {
        statement = "UPDATE post SET title=?, body=?, author=?, date=? WHERE id = ?"

    }

    insert, err := con.Prepare(statement)

    if err != nil {
        log.Printf("Error preparing statement: %v", err)
        return
    }

    p.PostDate = time.Now()

    if !p.isSaved {
        res, err = insert.Exec(p.Title, p.Body, p.Author, p.PostDate.Format(time_layout))
    } else {
        res, err = insert.Exec(p.Title, p.Body, p.Author, p.PostDate.Format(time_layout), p.ID)
    }

    if err != nil {
        fmt.Printf("Error inserting post into database: %v", err)
        return
    }

    lastID, err := res.LastInsertId()

    if err != nil {
        log.Printf("Error obtaining ID for inserted row: %v", err)
        return
    }

    p.ID = int(lastID)
    p.isSaved = true
}

func (p Post) String() string {
    return fmt.Sprintf("<h4>%v</h4><p>%v</p><small>By %v, %v</small>", p.Title, p.Body, p.Author, p.PostDate.Format(human_time_layout))
}

func GetPost(con *sql.DB, id int) Post {
    var err error
    row := con.QueryRow("SELECT title, body, author, date FROM post WHERE id = ?", id)

    if err != nil {
        log.Printf("Error preparing fetch statement: %v", err)
    }

    var date string
    post := new(Post)
    err = row.Scan(&post.Title, &post.Body, &post.Author, date)

    post.PostDate, err = time.Parse(time_layout, date)
    post.isSaved = true

    if err != nil {
        log.Printf("Error parsing date: %v \n - got %v", err, date)
    }

    return *post
}

func checkIndex(index []int, idx int) bool {
    result := false

    for count := 0; count < len(index); count++ {
        if index[count] == idx {
            result = true
            break
        }
    }

    return result
}

func GetPostList(con *sql.DB, position, limit int) []Post {
    var err error
    conn := *con
    posts := make([]Post, 0)

    rows, err := conn.Query("SELECT id, title, body, author, date FROM post LIMIT ?,?", position, limit)

    if err != nil {
        log.Print(err)
    }

    defer rows.Close()

    var post *Post
    var date string
    index := make([]int, 0)

    for rows.Next() {
        post = new(Post)
        err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Author, &date)

        post.PostDate, err = time.Parse(time_layout, date)
        post.isSaved = true

        if err != nil {
            log.Printf("Error parsing date: %v \n - got %v", err, date)
        }

        if !checkIndex(index, post.ID) {
            posts = append(posts, *post)
            index = append(index, post.ID)
        }
    }

    err = rows.Err()

    if err != nil {
        log.Print(err)
    }

    return posts
}
