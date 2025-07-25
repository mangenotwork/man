package main

import (
	"fmt"
	"regexp"
)

func main() {
	case1()
	case2()
	case3()
	case4()
}

// 在 Go 中使用蕴含逻辑的本质是用简洁的逻辑表达 “约束关系”，尤其适合：
// 1. 当条件 P 不成立时，规则自动生效（无需额外处理）； 2. 当条件 P 成立时，强制要求 Q 必须成立，否则触发异常。
// 这种模式能让代码更简洁、逻辑更清晰，避免冗余的嵌套判断（如 if P { if !Q { 报错 } }）。

// 1. 接口参数校验
// 例如：“如果用户提交了邮箱（P），则邮箱必须符合格式（Q）”。
func case1() {
	// 场景：用户提交的表单中，若填写了邮箱，则必须合法
	userInput := struct {
		Email string // 可能为空（未填写）
	}{
		Email: "invalid-email", // 填写了但格式错误
	}

	// P：用户填写了邮箱（Email 不为空）
	p := userInput.Email != ""
	// Q：邮箱格式合法
	q := isEmailValid(userInput.Email)

	// 蕴含逻辑：如果填写了邮箱（P），则必须合法（Q）
	isValid := !p || q

	if !isValid {
		fmt.Println("错误：邮箱格式无效")
	} else {
		fmt.Println("参数校验通过")
	}
}

// 验证邮箱格式
func isEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// 2. 权限控制 例如：“如果操作是删除（P），则用户必须是管理员（Q）”。
func case2() {
	// 场景：删除操作需要管理员权限
	action := "delete" // 操作类型：delete / view / edit 等
	isAdmin := false   // 当前用户是否为管理员
	// P：操作是删除
	p := action == "delete" // Q：用户是管理员
	q := isAdmin
	// 蕴含逻辑：如果是删除操作（P），则必须是管理员（Q）
	hasPermission := !p || q
	if hasPermission {
		fmt.Println("允许操作")
	} else {
		fmt.Println("权限不足：删除操作需要管理员权限")
	}
}

// 3. 避免空指针异常
// 例如：“如果对象不为 nil（P），则其属性必须满足某个条件（Q）”。
func case3() {

	type User struct {
		Name string
	}

	// 场景：如果用户对象存在（非 nil），则用户名不能为空
	var user *User = &User{Name: ""} // 用户存在，但姓名为空

	// P：用户对象非 nil
	p := user != nil
	// Q：用户姓名不为空
	q := user.Name != ""

	// 蕴含逻辑：如果用户存在（P），则姓名必须非空（Q）
	isUserValid := !p || q

	if isUserValid {
		fmt.Println("用户信息合法")
	} else {
		fmt.Println("错误：用户存在但姓名为空")
	}

}

// 异或的核心价值在于高效处理 “差异” 和 “可逆操作”，尤其在无需复杂逻辑的场景中，能以极简的代码实现功能，兼具性能优势。实际开发中，需根据具体需求（如加密、算法优化、状态控制等）灵活运用。

func case4() {
	a, b := true, false

	// 异或逻辑：a 和 b 恰好有一个为 true
	result := (a || b) && !(a && b)

	fmt.Println(result) // 输出：true

	// 布尔值异或（通过逻辑表达式）
	boolXOR := func(a, b bool) bool {
		return (a || b) && !(a && b)
	}

	fmt.Println(boolXOR(true, false))  // 输出：true
	fmt.Println(boolXOR(true, true))   // 输出：false
	fmt.Println(boolXOR(false, false)) // 输出：false

	// 整数按位异或
	x, y := 5, 3       // 5 是 101，3 是 011（二进制）
	fmt.Println(x ^ y) // 输出：6（110，二进制）

	// 应用：切换状态（true <-> false）
	status := true
	status = boolXOR(status, true) // 等价于 status = !status
	fmt.Println(status)            // 输出：false

	// 切换开关状态（true <-> false）
	status2 := true
	status2 = status2 != true // 常规写法
	// 等价于异或逻辑：status = (status || true) && !(status && true)
	fmt.Println(status2)
}

// 反蕴含在 Go 中的应用本质是校验 “结果状态” 对 “前提条件” 的依赖性，即：
//
//当某个结果（Q）成立时，必须确保其依赖的前提（P）已满足；
//若结果（Q）不成立，则前提（P）是否满足不影响逻辑的合法性。
//
//这种逻辑能简化状态依赖、权限校验等场景的代码，避免冗余的条件嵌套，使约束关系更清晰。

// 1. 数据依赖性校验 例如：“如果订单状态为‘已支付’（Q），则必须存在支付记录（P）”。
func case5() {
	type Order struct {
		Status    string // "unpaid" 或 "paid"
		PaymentID string // 支付记录ID，未支付时为空
	}

	// 场景：已支付的订单必须有支付记录
	order := Order{
		Status:    "paid", // Q：订单状态为已支付
		PaymentID: "",     // P：支付记录是否存在（非空即为存在）
	}

	q := order.Status == "paid" // Q：订单已支付
	p := order.PaymentID != ""  // P：存在支付记录

	// 反蕴含逻辑：如果订单已支付（Q），则必须有支付记录（P）
	isValid := !q || p

	if !isValid {
		fmt.Println("错误：已支付订单缺少支付记录")
	} else {
		fmt.Println("订单状态校验通过")
	}
}

// 2. 权限与操作的反向约束 例如：“如果用户能查看敏感数据（Q），则必须是高级会员（P）”。
func case6() {
	// 场景：能查看敏感数据的用户必须是高级会员
	canViewSensitive := true // Q：用户有权限查看敏感数据
	isPremiumMember := false // P：用户是否为高级会员
	q := canViewSensitive
	p := isPremiumMember
	// 反蕴含逻辑：如果能查看敏感数据（Q），则必须是高级会员（P）
	hasValidPermission := !q || p
	if !hasValidPermission {
		fmt.Println("错误：非高级会员无权查看敏感数据")
	} else {
		fmt.Println("权限校验通过")
	}
}

// 3. 状态机合法性检查
// 例如：“如果设备处于‘运行中’状态（Q），则必须已完成初始化（P）”。
func case7() {
	type Device struct {
		State       string // "idle" / "running" / "stopped"
		Initialized bool   // 是否已完成初始化
	}

	// 场景：运行中的设备必须已初始化
	device := Device{
		State:       "running", // Q：设备运行中
		Initialized: false,     // P：是否已初始化
	}

	q := device.State == "running" // Q：设备运行中
	p := device.Initialized        // P：已初始化

	// 反蕴含逻辑：如果设备运行中（Q），则必须已初始化（P）
	isStateValid := !q || p

	if !isStateValid {
		fmt.Println("错误：设备未初始化，无法进入运行状态")
	} else {
		fmt.Println("设备状态合法")
	}
}
