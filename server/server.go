package server

import (
	"goly/model"
	"goly/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func redirect(ctx *fiber.Ctx) error {
	golyrl := ctx.Params("redirect")
	goly, err := model.FindByGolyUrl(golyrl)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not find redirect",
			"error":   err,
		})
	}

	// grab any stats you want...
	goly.Clicked += 1
	_, err = model.UpdateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not update goly",
			"error":   err,
		})
	}

	return ctx.Redirect(goly.Redirect, fiber.StatusTemporaryRedirect)
}

func healthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Server is running",
	})
}

func getAllGolies(ctx *fiber.Ctx) error {
	golies, err := model.GetAllGolies()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get all redirects",
			"error":   err,
		})
	}
	// Return the array of JSON objects
	return ctx.Status(fiber.StatusOK).JSON(golies)
}

func getGoly(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not parse id",
			"error":   err,
		})
	}

	goly, err := model.GetGoly(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not get goly",
			"error":   err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)
}

func createGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")

	var goly model.Goly
	err := ctx.BodyParser(&goly)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not parse JSON",
			"error":   err,
		})
	}

	if goly.Random {
		goly.Goly = utils.RandomURL(8)
	}

	newGoly, err := model.CreateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not create goly",
			"error":   err,
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(newGoly)
}

func updateGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")

	var goly model.Goly
	err := ctx.BodyParser(&goly)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not parse JSON",
			"error":   err,
		})
	}

	updatedGoly, err := model.UpdateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not update goly",
			"error":   err,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(updatedGoly)
}

func deleteGoly(ctx *fiber.Ctx) error {
	var err error
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not parse id",
			"error":   err,
		})
	}

	err = model.DeleteGoly(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not delete goly",
			"error":   err,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Goly deleted",
	})
}

func SetupAndListen() {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Get("/", healthCheck)
	router.Get("/r/:redirect", redirect)
	router.Get("/goly", getAllGolies)
	router.Get("/goly/:id", getGoly)
	router.Post("/goly", createGoly)
	router.Patch("/goly", updateGoly)
	router.Delete("/goly/:id", deleteGoly)

	router.Listen(":3001")
}
