package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/reztimo/fiber-mongodb-hrms/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeController struct {
	Db *mongo.Database
}

func NewEmployeeController(db *mongo.Database) *EmployeeController {
	return &EmployeeController{
		Db: db,
	}
}

func (ec *EmployeeController) GetEmployees(c *fiber.Ctx) error {
	query := bson.D{{}}

	cursor, err := ec.Db.Collection("employees").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var employees []models.Employee
	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employees)
}

func (ec *EmployeeController) CreateEmployee(c *fiber.Ctx) error {
	collection := ec.Db.Collection("employees")

	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	employee.ID = primitive.NewObjectID()

	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdEmployee := &models.Employee{}
	if err := createdRecord.Decode(createdEmployee); err != nil {
		log.Println("Error decoding created employee:", err)
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(createdEmployee)
}

func (ec *EmployeeController) UpdateEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).SendString("Invalid ID format")
	}

	employee := new(models.Employee)
	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: employeeID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}

	result := ec.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(404).SendString("Employee not found")
		}
		return c.Status(500).SendString(result.Err().Error())
	}

	employee.ID = employeeID

	return c.Status(200).JSON(employee)
}

func (ec *EmployeeController) DeleteEmployee(c *fiber.Ctx) error {
	employeeID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	query := bson.D{{Key: "_id", Value: employeeID}}
	result, err := ec.Db.Collection("employees").DeleteOne(c.Context(), query)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Employee not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Record deleted"})
}
