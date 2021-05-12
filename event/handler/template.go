package handler

import (
	"bytes"
	"text/template"
)

var (
	personalInfo, _ = template.New("personalInfo").Parse(`这是你的契约喵~
昵称: {{.UserInfo.Nickname}}
通行证: {{.UserInfo.UID}}`)
)

func useTemplate(t *template.Template, obj interface{}) string {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, obj)
	if err != nil {
		return "铲屎官好像写错东西了喵~\n" + err.Error()
	}
	return buf.String()
}

func runtimeErr(err error) string {
	return "好像出错了喵~\n" + err.Error()
}

func applicationErr(err error) string {
	return "啊咧~\n" + err.Error()
}
