package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

type ActionCardBtn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type ActionCardMsg struct {
	Title          string           `json:"title"`
	Text           string           `json:"text"`
	Btns           *[]ActionCardBtn `json:"btns"`
	BtnOrientation string           `json:"btnOrientation"`
	HideAvatar     string           `json:"hideAvatar"`
}

type Msg struct {
	MsgType    string      `json:"msgtype"`
	ActionCard interface{} `json:"actionCard"`
}

func getTitle(model hookModel) string {
	switch model.ObjectKind {
	case "build":
		return fmt.Sprintf("Pipeline: [%s] %s", model.Project.Name, model.ObjectAttributes.Status)
	case "pipeline":
		return fmt.Sprintf("Job: [%s] %s @ %s ( %s )", model.ProjectName, model.BuildName, model.BuildStage, model.BuildStatus)
	default:
		return "Unsupported message kind"
	}
}

func getText(model hookModel) string {
	t := getHookTemplate(model)
	if t == nil {
		t = getDefaultTemplate(model)
		if t == nil {
			return "Notfound template: " + model.ObjectAttributes.Status
		}
	}
	var tpl bytes.Buffer
	err := t.Execute(&tpl, model)
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func getDefaultTemplate(m hookModel) *template.Template {
	return getTemplate(m.ObjectKind + ".default.mdt")
}

var templates = make(map[string]*template.Template)

func getHookTemplate(model hookModel) *template.Template {
	return getTemplate(getTemplateName(model))
}

func getTemplate(templateName string) *template.Template {
	if t, ok := templates[templateName]; ok {
		return t
	}
	info, err := os.Stat(templateName)
	if os.IsNotExist(err) || info.IsDir() {
		return nil
	}
	t, err := template.New(templateName).ParseFiles("./" + templateName)
	if err != nil {
		panic(err)
	}
	templates[templateName] = t
	return t
}

func getTemplateName(model hookModel) string {
	switch model.ObjectKind {
	case "build":
		return model.BuildStage + "." + model.ObjectKind + ".mdt"
	default:
		return model.ObjectAttributes.Status + ".mdt"
	}
}

func getMsgBody(model hookModel) string {
	var btns []ActionCardBtn
	actionCard := ActionCardMsg{
		Title:          getTitle(model),
		Text:           getText(model),
		Btns:           &btns,
		BtnOrientation: "0",
		HideAvatar:     "0",
	}
	msg := Msg{
		MsgType:    "actionCard",
		ActionCard: &actionCard,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(b)
}
