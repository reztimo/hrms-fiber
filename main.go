// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/reztimo/fiber-mongodb-hrms/controllers"
	"github.com/reztimo/fiber-mongodb-hrms/db"
)

func main() {
	app := fiber.New()

	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	employeeController := controllers.NewEmployeeController(db.Mg.Db)

	app.Get("/employee", employeeController.GetEmployees)
	app.Post("/employee", employeeController.CreateEmployee)
	app.Put("/employee/:id", employeeController.UpdateEmployee)
	app.Delete("/employee/:id", employeeController.DeleteEmployee)

	app.Listen(":5151")
}
