package main

import (
	"encoding/json"
	"fmt"
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

func getPipelineLink(model hookModel) string {
	return fmt.Sprintf("%s/pipelines/%d", model.Project.WebUrl, model.ObjectAttributes.Id)
}

func getTitle(model hookModel) string {
	return fmt.Sprintf("Run [%s] pipeline: %s", model.Project.Name, model.ObjectAttributes.Status)
}

func getText(model hookModel) string {
	content := fmt.Sprintf(`The pipeline is triggered by the commit **\"%s\"** pushed by %s(%s)`,
		model.Commit.Message, model.Commit.Author.Name, model.Commit.Author.Email)
	return content
}

func getMsg(model hookModel) string {
	btns := []ActionCardBtn{
		{
			Title:     "查看详情",
			ActionURL: getPipelineLink(model),
		},
	}
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
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
