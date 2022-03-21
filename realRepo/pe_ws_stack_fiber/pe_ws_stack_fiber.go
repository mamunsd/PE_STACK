package pe_ws_stack_fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mamunsd/PE_STACK/pe_file_ops"
)

var homeCache string

func PeFwebServerSample() {
	homeCache = pe_file_ops.GetFileContentAsString("./statics/page0.html")
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1"})
	app.Static("/", "./static")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH",
		// AllowMethods: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/page0.html")
	})

	regProdPostRoutes(app)
	// regSampleRoutesSysExe(app)
	// regSamplePageRoutes(app)

	app.Listen("0.0.0.0:9928")
}
