package main

import (
	"be-evaporator/database"
	_ "be-evaporator/docs"
	"be-evaporator/routes"
	"be-evaporator/utils"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// @title REST API for Evaporator Datastore
// @version 2024.07.02.1
// @description REST API for Evaporator Datastore

// @contact.name Pai Joe
// @contact.email pai.joe@wedibojone.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @Host eva-admin.msvc.app
// @BasePath /rest/api
func main() {

	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDatabase()
	defer func() {
		dbInstance, _ := database.DB.DB()
		_ = dbInstance.Close()
	}()

	limitUploadSize, _ := strconv.Atoi(os.Getenv("LIMIT_UPLOAD_SIZE"))

	app := fiber.New(fiber.Config{
		JSONDecoder:    json.Unmarshal,
		JSONEncoder:    json.Marshal,
		BodyLimit:      limitUploadSize * 1024 * 1024,
		ReadBufferSize: 10 * 1024 * 1024,
		// Override default error handler
		ErrorHandler: processError,
	})

	// setup middleware
	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(etag.New())
	app.Use(logger.New())

	// Add recovery middleware
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize routes
	routes.Setup(app)

	err = app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
		return
	}
}

func processUtilError(ctx *fiber.Ctx, err *utils.Error) error {
	return ctx.Status(err.Status).JSON(err)
}

func processFiberError(ctx *fiber.Ctx, err *fiber.Error) error {
	return ctx.Status(err.Code).JSON(utils.Error{Status: err.Code, Code: "internal-server", Message: err.Message})
}

func processDefaultError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(500).JSON(utils.Error{Status: 500, Code: "internal-server", Message: err.Error()})
}

func processError(ctx *fiber.Ctx, err error) error {
	var e *utils.Error
	if errors.As(err, &e) {
		return processUtilError(ctx, e)
	}

	var f *fiber.Error
	if errors.As(err, &f) {
		return processFiberError(ctx, f)
	}

	return processDefaultError(ctx, err)
}
