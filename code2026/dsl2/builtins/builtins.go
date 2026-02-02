package builtins

/*

内置函数注册器

*/

import (
	"dsl2/interpreter"
	"fmt"
	"log"
	"strings"
	"time"
)

// 注册所有内置函数
func RegisterBuiltins(interp *interpreter.Interpreter) {
	registerMath(interp)
	registerString(interp)
	registerTime(interp)
	registerUtils(interp)
	registerCore(interp)
}

func registerMath(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("abs", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("abs() 需要一个参数")
		}

		switch v := args[0].(type) {
		case int64:
			if v < 0 {
				return -v, nil
			}
			return v, nil
		case float64:
			if v < 0 {
				return -v, nil
			}
			return v, nil
		default:
			return nil, fmt.Errorf("abs() 不支持的类型: %T", args[0])
		}
	})

	interp.Global().SetFunc("max", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("max() 需要至少2个参数")
		}

		// 实现max函数
		return nil, nil
	})

	interp.Global().SetFunc("min", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("min() 需要至少2个参数")
		}

		// 实现min函数
		return nil, nil
	})
}

func registerString(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("upper", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("upper() 需要一个参数")
		}

		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("upper() 需要字符串参数")
		}

		return strings.ToUpper(s), nil
	})

	interp.Global().SetFunc("lower", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("lower() 需要一个参数")
		}

		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("lower() 需要字符串参数")
		}

		return strings.ToLower(s), nil
	})

	interp.Global().SetFunc("trim", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("trim() 需要一个参数")
		}

		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("trim() 需要字符串参数")
		}

		return strings.TrimSpace(s), nil
	})

	interp.Global().SetFunc("split", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("split() 需要2个参数")
		}

		s, ok1 := args[0].(string)
		sep, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("split() 需要字符串参数")
		}

		parts := strings.Split(s, sep)
		result := make([]interpreter.Value, len(parts))
		for i, part := range parts {
			result[i] = part
		}
		return result, nil
	})
}

func registerTime(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("now", func(args []interpreter.Value) (interpreter.Value, error) {
		return time.Now().Unix(), nil
	})

	interp.Global().SetFunc("sleep", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("sleep() 需要一个参数")
		}

		var ms int64
		switch v := args[0].(type) {
		case int64:
			ms = v
		case float64:
			ms = int64(v)
		default:
			return nil, fmt.Errorf("sleep() 需要数字参数")
		}

		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil, nil
	})
}

func registerUtils(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("type_of", func(args []interpreter.Value) (interpreter.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("type_of() 需要一个参数")
		}

		switch args[0].(type) {
		case int64:
			return "int", nil
		case float64:
			return "float", nil
		case string:
			return "string", nil
		case bool:
			return "bool", nil
		default:
			return "unknown", nil
		}
	})

	interp.Global().SetFunc("exit", func(args []interpreter.Value) (interpreter.Value, error) {
		code := 0
		if len(args) > 0 {
			if v, ok := args[0].(int64); ok {
				code = int(v)
			}
		}
		// 这里应该优雅退出，但为了简单示例，我们只设置一个标志
		panic(fmt.Sprintf("exit(%d)", code))
	})
}

func registerCore(interp *interpreter.Interpreter) {
	interp.Global().SetFunc("chrome", func(args []interpreter.Value) (interpreter.Value, error) {
		log.Println("执行 chrome 的操作，参数是 ", args, len(args))
		return nil, nil
	})
}
