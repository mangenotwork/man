package main

import "github.com/mangenotwork/gathertool"

func main() {
	test := gathertool.NewTestUrl("https://10.0.40.3/webapi/activity/info", "POST", 50000, 5000)
	test.JsonData = "{\n    \"activity_id\": 299\n}"
	test.Run()
}
