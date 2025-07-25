package main

import (
//    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
//    "github.com/redis/go-redis/v9"
//    "encoding/json"
//    "log"
)


func setRoutes(app *fiber.App) {

    app.Post("/signup", func(c *fiber.Ctx) error {

        // Get the name from the form
        username := c.FormValue("nusername")
        password := c.FormValue("npassword")

        // Print it server side
        fmt.Printf("signing up %s with password %s\n", username, password)

        
        if (key_exists(username)) {
            fmt.Printf("user exists already")
            return c.Render("signupPage", fiber.Map{
                "warning": "User already exists! Log in instead!",
                "Name": "User",
            })
        } else {

			//~~~~~~~~~~~~~~~~~~~successful sign up~~~~~~~~~~~~~~~~~~~~~~~~~~
        
			//=================================================================
            //add the name to the database
            if err := create_new_user(username, password); err != nil {
                fmt.Printf("adding user failed")
                panic(err)
            }
		    //create a session
			sess, err := store.Get(c)
			if err != nil {
				return err
			}
			
			//store the username in session
			sess.Set("username", username)
			
			//save the session
			if err := sess.Save(); err != nil {
				return err
			}
			
			//redirect to profile
			return c.Redirect("/userProfile")


			//==================================================
        }

    })


    app.Post("/login", func(c *fiber.Ctx) error {

        fmt.Println("/login route triggered")

        username := c.FormValue("username")
        password := c.FormValue("password")

        fmt.Printf("username=%s, password=%s\n", username, password)

        if !key_exists(username) {

            fmt.Println("DEBUG: User doesn't exist")

            return c.Render("loginPage", fiber.Map{
                "warning": "Username not found! Sign up!",
                "Name": "User",
            })

        } else {
            value, _ := GetUserData(username, "password");

            if (value != password){
                fmt.Printf("wrong password")

                return c.Render("loginPage", fiber.Map{
                    "warning": "Wrong password, check your spelling",
                    "Name": "User",

                })
            } else {

				//======================successful log in=================================
                fmt.Println("Login successful")
                
				//create a session
				sess, err := store.Get(c)
				if err != nil {
					return err
				}
				
				//store the username in session
				sess.Set("username", username)
				
				//save the session
				if err := sess.Save(); err != nil {
					return err
				}
				
				//redirect to profile
				return c.Redirect("/userProfile")


				//======================================================================

            }
        }
        

    })


    app.Post("/userInfo", func(c *fiber.Ctx) error {
        username := c.FormValue("username") 
        password := c.FormValue("password") 
        name := c.FormValue("name")
        age := c.FormValue("age")
        color := c.FormValue("color")

        if err := SetUserData(username, "name", name); err != nil {
            fmt.Printf("adding name failed")
            panic(err)
        }
        if err := SetUserData(username, "age", age); err != nil {
            fmt.Printf("adding age failed")
            panic(err)
        }
        if err := SetUserData(username, "color", color); err != nil {
            fmt.Printf("adding color failed")
            panic(err)
        }

        return c.Render("logged-in", fiber.Map{
            "username": username,
            "password": password,
            "name": name,
            "age": age,
            "color": color,
        })
    })

	app.Get("/loginPage", func(c *fiber.Ctx) error {
		return c.Render("loginPage", nil)
	})
	
	app.Get("/index", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
	
	
	app.Get("/signupPage", func(c *fiber.Ctx) error {
		return c.Render("signupPage", nil)
	})

	app.Get("/userProfile", func(c *fiber.Ctx) error {

		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		    
		//get username from session
		username := sess.Get("username")
		
		if username == nil {
			return c.Redirect("/") //if not logged in, go to homepage
		}



		var userName, _ = GetUserData(username.(string), "name")
		var userAge, _ = GetUserData(username.(string), "age")
		var userColor, _ = GetUserData(username.(string), "color")

		return c.Render("userProfile", fiber.Map{
			"username": username.(string),
			"name": userName,
			"age": userAge,
			"color": userColor,
		})



	})

    fmt.Println("Finished route setup")

}
