package main

import (
//    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "strconv"  
//    "github.com/redis/go-redis/v9"
//    "encoding/json"
//    "log"
    "time"
    "strings"
)


func setRoutes(app *fiber.App) {


    //====================[FORMS]========================

    //what happens when user hits sign up button:
    app.Post("/signup", func(c *fiber.Ctx) error {

        //get the user and pass from the form
        username := c.FormValue("username")
        password := c.FormValue("password")
        name := c.FormValue("name")
        age := c.FormValue("age")
        email := c.FormValue("email")

        //debug console print
        fmt.Printf("signing up %s with password %s\n", username, password)

        
        if (key_exists(username)) {

            //prints error message to console
            fmt.Printf("user exists already")

            //redisplays site with new warning
            return c.Render("signupPage", fiber.Map{
                "AddError": "User already exists! Log in instead!",
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
            if err := SetUserData(username, "name", name); err != nil {
                fmt.Printf("adding name failed")
                panic(err)
            }
            if err := SetUserData(username, "age", age); err != nil {
                fmt.Printf("adding age failed")
                panic(err)
            }
            if err := SetUserData(username, "email", email); err != nil {
                fmt.Printf("adding color failed")
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

    //what happens when user hits log in button
    app.Post("/login", func(c *fiber.Ctx) error {

        fmt.Println("/login route triggered")

        username := c.FormValue("username")
        password := c.FormValue("password")

        fmt.Printf("username=%s, password=%s\n", username, password)

        if !key_exists(username) {

            fmt.Println("DEBUG: User doesn't exist")

            return c.Render("loginPage", fiber.Map{
                "AddError": "Username not found! Sign up!",
                "Name": "User",
            })

        } else {
            value, _ := GetUserData(username, "password");

            if (value != password){
                fmt.Printf("wrong password")

                return c.Render("loginPage", fiber.Map{
                    "AddError": "Wrong password, check your spelling",
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

    //what happens when user submits their info
    app.Post("/userInfo", func(c *fiber.Ctx) error {

        //collect the username and the form values
        username := getUser(c)
        name := c.FormValue("name")
        age := c.FormValue("age")
        color := c.FormValue("color")
        email := c.FormValue("email")

        //set the data into the database or display error on console
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
        if err := SetUserData(username, "email", email); err != nil {
            fmt.Printf("adding color failed")
            panic(err)
        }

        //rerender page
        return c.Redirect("/userProfile")
    })

    //what happens when user adds workout
    app.Post("/addWorkout/add", func(c *fiber.Ctx) error {
        //get all data
        username := getUser(c)

        if username == "" {
            return c.Redirect("/") // if not logged in, go to homepage
        }
    
        exercise := strings.ToLower(c.FormValue("exercise"))
        sets, _ := strconv.Atoi(c.FormValue("sets"))
        reps, _ := strconv.Atoi(c.FormValue("reps"))
        weight, _ := strconv.Atoi(c.FormValue("weight"))
        date := c.FormValue("date") // Get date from form
        
        // If no date specified, use today's date
        if date == "" {
            currentTime := time.Now()
            date = currentTime.Format("2006-01-02")
        }
        
        //add workout and collect error
        err := AddWorkout(username, date, exercise, [3]int{sets, reps, weight})
        
        //rerender page with new error if error happened
        if err != nil {
            return c.Render("addWorkout", fiber.Map{
                "AddError": err.Error(),
                "Username": username,
                "Date":     date,
            })
        }
        
        //rerender with success message
        return c.Render("addWorkout", fiber.Map{
            "Username": username,
            "Date":     date,
            "Success":  "Workout added successfully!",
        })
    })

    
    app.Get("/displayWOs", func(c *fiber.Ctx) error {
        
        //get username from session and turn to string
        username := getUser(c)

        if username == "" {
            return c.Redirect("/") //if not logged in, go to homepage
        }

        date := c.FormValue("date")
        var origin string 
        origin = c.FormValue("origin")
        
        var workout map[string]interface{}
        var err error
        
        //if date isnt empty get the workouts from the day and store it in workout
        if date != "" {
            workout, err = GetWorkout(username, date)

        }
        
        if (origin == "addWorkout") {
        //rerender website to have workout, or maybe error
            return c.Render("addWorkout", fiber.Map{
                "Username": username,
                "Date":     date,
                "Workout":  workout,
                "Error":    err,
            }) 
        } else {    
            
            return c.Render("calendar", fiber.Map{
                "Username": username,
                "Date":     date,
                "Workout":  workout,
                "Error":    err,
            })

        }


    })
    

    app.Post("/generateGraph", func(c *fiber.Ctx) error {
        username := getUser(c)
        if username == "" {
            return c.Redirect("/") //if not logged in, go to homepage
        }
    
        //get form data
        startDate := c.FormValue("startDate")
        endDate := c.FormValue("endDate")
        exercise := strings.ToLower(c.FormValue("workout"))
    
        //get all workouts for this exercise
        workouts, err := GetWorkoutsBetweenDates(username, startDate, endDate, exercise)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
    
        var dates []string
        var weights []int
        var reps []int
        var sets []int
    
        //loop through workouts and saave into new arrays
        for _, workout := range workouts {
            dates = append(dates, workout.Date)
            weights = append(weights, workout.Weight)
            reps = append(reps, workout.Reps)
            sets = append(sets, workout.Sets)
        }
    
        return c.JSON(fiber.Map{
            "exercise": exercise,
            "dates":    dates,
            "weights":  weights,
            "reps":     reps,
            "sets":     sets,
        })
    
    })

    app.Post("/logout", func(c *fiber.Ctx) error {

        username := getUser(c)
        if username == "" {
            return c.Redirect("/") //if not logged in, go to homepage
        }

        sess, err := store.Get(c)
        if err != nil {
            panic(err)
        }
        
        //destroy session
        if err := sess.Destroy(); err != nil {
            panic(err)
        }
        
        return c.Redirect("/")
    })

    //====================[PAGES]========================

    //what happens when user visits the main page
	app.Get("/index", func(c *fiber.Ctx) error {
        username := getUser(c)
        if username != "" {
            return c.Redirect("/home") //if logged in, go to logged in version
        }
		return c.Render("index", nil)
	})

    //what happens when user visits the logged in version of the main page
    app.Get("/home", func(c *fiber.Ctx) error {
        username := getUser(c)
        if username == "" {
            return c.Redirect("/") //if not logged in, go to not logged in version
        }
		return c.Render("home", nil)
	})

    //what happens when user visits the log in page
	app.Get("/loginPage", func(c *fiber.Ctx) error {

        //if user is logged in, bring them to their page
        username := getUser(c)
        if username != "" {
            return c.Redirect("/userProfile")
        }

        //else render log in page
		return c.Render("loginPage", nil)
	})
	
	//what happens when user visits the sign in page
	app.Get("/signupPage", func(c *fiber.Ctx) error {

        //if user is logged in, bring them to their page
        username := getUser(c)
        if username != "" {
            return c.Redirect("/userProfile") 
        }
        //else render sign up page
		return c.Render("signupPage", nil)
	})

    //what happens when a user visits their profile
	app.Get("/userProfile", func(c *fiber.Ctx) error {

        //get the username of the session, if no session, go back to home
        username := getUser(c)
        if username == "" {
            return c.Redirect("/") //if not logged in, go to homepage
        }

        //setting variables to the user data
        //we cant plug these directly into the rendering because the function
        //returns an err too, which we ignore using '_'
		var userName, _ = GetUserData(username, "name")
		var userAge, _ = GetUserData(username, "age")
		var userColor, _ = GetUserData(username, "color")
        var userSince, _ = GetUserData(username, "since")
        var userEmail, _ = GetUserData(username, "email")

        //render using user data
		return c.Render("userProfile", fiber.Map{
			"username": username,
			"name": userName,
			"age": userAge,
			"color": userColor,
            "since": userSince,
            "email": userEmail,
		})

	})

    //add workout page
    app.Get("/addWorkout", func(c *fiber.Ctx) error {
        username := getUser(c)
        if username == "" {
            return c.Redirect("/") // if not logged in, go to homepage
        }
    
        currentTime := time.Now()
        date := currentTime.Format("2006-01-02")
    
        return c.Render("addWorkout", fiber.Map{
            "Username": username,
            "Date":     date,
        })
    })

    //calendar page
    app.Get("/calendar", func(c *fiber.Ctx) error {
        return c.Render("calendar", nil)
    })

    //create graph page
    app.Get("/createGraph", func(c *fiber.Ctx) error {
        return c.Render("createGraph", nil)
    })

    //features page
    app.Get("/features", func(c *fiber.Ctx) error {
        return c.Render("features", nil)
    })

    //timer page
    app.Get("/timer", func(c *fiber.Ctx) error {
        return c.Render("timer", nil)
    })

    //reminders page
    app.Get("/reminders", func(c *fiber.Ctx) error {
        return c.Render("reminders", nil)
    })


    
//=================[TEST AREA]===================


//=========================

















    fmt.Println("Finished route setup")

}
