package pe_temp_funcs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mamunsd/PE_STACK/pe_mongo_db"
)

var homeCache string

func GetFileContentAsByte(filePath string) []byte {
	myFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	myByteVal, _ := ioutil.ReadAll(myFile)
	return myByteVal
}

func GetFileContentAsString(filePath string) string {
	myFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	myByteVal, _ := ioutil.ReadAll(myFile)
	return string(myByteVal)
}

func PeSysCmd(myCommand string) string {
	cmd := exec.Command("/bin/sh", "-c", myCommand)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	myOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(myOutput)
}

func PeSysCmdWait(myCommand string) string {
	myexecute := exec.Command("/bin/sh", "-c", myCommand+" &")
	myOutput, err := myexecute.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(myOutput)
}

func regSampleRoutesSysExe(app *fiber.App) {
	app.Get("/api/v1/system/state/pingHost", func(c *fiber.Ctx) error {
		return c.SendString(PeSysCmdWait("fping -q -a -s -c 70 -p 10 -t 10 192.168.56.100 192.168.56.120 192.168.56.121 192.168.56.122 103.99.250.145"))
	})
}

func regSamplePageRoutes(app *fiber.App) {
	page1String := GetFileContentAsString("spaUiFiles/page1.html")
	app.Get("/spaUi/page1", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(page1String)
	})
}

func regPostRoutes(app *fiber.App) {
	app.Post("/api/queryEcho", func(c *fiber.Ctx) error {
		Mreply := "আপাতত নাই"
		return c.SendString(Mreply)
	})

	app.Get("/api/ottojs", func(c *fiber.Ctx) error {
		// myReq := `{"q_collection":"short_messages","filter":{"_id": {"$gt": "akashetejototara"}, "USER_ID" : "rasel@gmail.com"},"qconfig":{"sort":{"_id":-1},"limit":75}}`
		// mReply := peMongo.GeneralQuery_withMap([]byte(myReq))
		return c.SendString("who")
	})

	app.Get("/api/gjson", func(c *fiber.Ctx) error {
		// myReq := `{"q_collection":"short_messages","filter":{"_id": {"$gt": "akashetejototara"}, "USER_ID" : "rasel@gmail.com"},"qconfig":{"sort":{"_id":-1},"limit":75}}`
		// mReply := peMongo.GeneralQuery_withMap([]byte(myReq))
		return c.SendString("")
	})

	app.Post("/api/generalQuery", func(c *fiber.Ctx) error {
		// myReq := `{"q_collection":"short_messages","filter":{"_id": {"$gt": "akashetejototara"}, "USER_ID" : "rasel@gmail.com"},"qconfig":{"sort":{"_id":-1},"limit":75}}`
		// mReply := peMongo.GeneralQuery_withMap([]byte(myReq))
		payload := c.Body()
		resString := pe_mongo_db.GenQueryMongo([]byte(payload))
		return c.SendString(resString)
	})
}

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

func PeFwebServerSample() {
	homeCache = GetFileContentAsString("./static/page0.html")
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
