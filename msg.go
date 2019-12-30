package main

import (
	"encoding/json"
	"fmt"
)

const (
	statusTitleTemplate = "%s <span style='color: %s'>%s</span>"

	statusSuccess = "success"
	statusPending = "pending"
	statusRunning = "running"
	statusFailed  = "failed"

	iconSuccess = "👍"
	iconPending = "⌛️"
	iconRunning = "🕘"
	iconFailed  = "💥"

	titleSuccess = "Success"
	titlePending = "Pending"
	titleRunning = "Running"
	titleFailed  = "Failed"

	colorSuccess = "Green"
	colorPending = "Blue"
	colorRunning = "Orange"
	colorFailed  = "Red"
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

func getStatusTitle(model hookModel) string {
	switch model.ObjectAttributes.Status {
	case statusSuccess:
		return fmt.Sprintf(statusTitleTemplate, iconSuccess, colorSuccess, titleSuccess)
	case statusPending:
		return fmt.Sprintf(statusTitleTemplate, iconPending, colorPending, titlePending)
	case statusRunning:
		return fmt.Sprintf(statusTitleTemplate, iconRunning, colorRunning, titleRunning)
	case statusFailed:
		return fmt.Sprintf(statusTitleTemplate, iconFailed, colorFailed, titleFailed)
	}
	return "???"
}

func getText(model hookModel) string {
	content := fmt.Sprintf(`### Run pipeline:

Project: **%s**

Ref: %s

Status: %s

Commit: %s

Author: %s

[查看详情](%s)`,
		model.Project.Name, model.ObjectAttributes.Ref, getStatusTitle(model),
		model.Commit.Message, model.Commit.Author.Name,
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
