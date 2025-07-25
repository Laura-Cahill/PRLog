package main

import (
    "context"
//    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
//    "encoding/json"
    "github.com/gofiber/template/html/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
    "log"
)

//user info struct to put in database
type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Name     string `json:"name"`
    Age      int    `json:"age"`
    Color    string `json:"color"`
}

var (
    //saves these where all go code can see it
    ctx = context.Background()
    rdb *redis.Client 
    app = fiber.New()
    store = session.New()
)

func main() {
    //init redis and connect to it
    initRedis()
	defer rdb.Close() //closes when we're done
    
    //initialize template engine
    engine := html.New("./html", ".html") // Changed from "./views" to "./html"
    
    //create Fiber instance
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    //serve static files once we need them
    app.Static("/css", "./css")
    app.Static("/js", "./js")

    //sets all the routes we'll need later
    setRoutes(app)


    //renders our homepage
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("index", fiber.Map{
            "Name": "User",
        })
    })

    //listens for users accessing website, logs if doesnt work
    log.Fatal(app.Listen("0.0.0.0:3000"))

}