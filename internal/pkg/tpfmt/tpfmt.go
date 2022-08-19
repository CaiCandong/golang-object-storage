package tpfmt

import (
	"strings"
	"text/template"
)

type FormatTp struct {
	tp *template.Template
}

// Exec 传入map填充预定的模板
func (f FormatTp) Exec(args map[string]interface{}) string {
	s := new(strings.Builder)
	err := f.tp.Execute(s, args)
	if err != nil {
		// 放心吧，这里不可能触发的，除非手贱:)
		panic(err)
	}
	return s.String()
}

/*
Format 自定义命名format，严格按照 {{.CUSTOMNAME}} 作为预定参数，不要写任何其它的template语法

usage:

	s = Format("{{.name}} hello.").Exec(map[string]interface{}{
	    "name": "superpig",
	}) // s: superpig hello.
*/
func Format(fmt string) FormatTp {
	var err error
	temp, err := template.New("").Parse(fmt)
	if err != nil {
		// 放心吧，这里不可能触发的，除非手贱:)
		panic(err)
	}
	return FormatTp{tp: temp}
}
