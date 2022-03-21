package pe_ws_stack_fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mamunsd/PE_STACK/pe_mongo_db"
)

func regProdPostRoutes(app *fiber.App) {
	app.Post("/api/sms/genq", func(c *fiber.Ctx) error {
		// myReq := `{"q_collection":"short_messages","filter":{"_id": {"$gt": "akashetejototara"}, "USER_ID" : "rasel@gmail.com"},"qconfig":{"sort":{"_id":-1},"limit":75}}`
		// mReply := peMongo.GeneralQuery_withMap([]byte(myReq))
		payload := c.Body()
		resString := pe_mongo_db.ShrtMsgGenQ(string(payload))
		return c.SendString(resString)
	})

	app.Post("/api/subs/new", func(c *fiber.Ctx) error {
		payload := c.Body()
		collName := "subscribers"
		pe_mongo_db.InsertOne(collName, string(payload))

		return c.SendString("{}")
	})
}
