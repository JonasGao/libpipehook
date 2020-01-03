package main

import (
	"encoding/json"
	"fmt"
)

const (
	statusSuccess = "success"
	statusPending = "pending"
	statusRunning = "running"
	statusFailed  = "failed"

	iconSuccess = "👍"
	iconPending = "⌛"
	iconRunning = "🕘"
	iconFailed  = "❌"

	titleSuccess = "Success"
	titlePending = "Pending"
	titleRunning = "Running"
	titleFailed  = "Failed"
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

func getStatusIcon(model hookModel) string {
	switch model.ObjectAttributes.Status {
	case statusSuccess:
		return iconSuccess
	case statusPending:
		return iconPending
	case statusRunning:
		return iconRunning
	case statusFailed:
		return iconFailed
	}
	return "???"
}

func getStatusTitle(model hookModel) string {
	switch model.ObjectAttributes.Status {
	case statusSuccess:
		return titleSuccess
	case statusPending:
		return titlePending
	case statusRunning:
		return titleRunning
	case statusFailed:
		return titleFailed
	}
	return "???"
}

func getText(model hookModel) string {
	content := fmt.Sprintf(`## Pipeline: %s %s

Project: **%s**

Commit: %s

Author: %s(%s)

[查看详情](%s)`,
		getStatusIcon(model), getStatusTitle(model), model.Project.Name,
		model.Commit.Message, model.Commit.Author.Name, model.Commit.Author.Email,
		getPipelineLink(model))
	return content
}

func getMsg(model hookModel) string {
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
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
