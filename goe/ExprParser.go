package goe

import (
	"bytes"
	"fmt"
	"text/template"
)

type Expr string //表达式类型

// ExecExpr 执行表达式，使用模板引擎 临时方法后期需要修改
func ExecExpr(expr Expr, data map[string]interface{}) (string, error) {
	tpl := template.New("expr").Funcs(map[string]interface{}{
		"echo": func(params ...interface{}) interface{} {
			return fmt.Sprintf("echo:%v", params[0])
		},
	})
	t, err := tpl.Parse(fmt.Sprintf("{{%s}}", expr))
	if err != nil {
		return "", err
	}
	var buf = &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
