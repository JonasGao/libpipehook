package main

import (
	"fmt"
	"testing"
)

func TestGetText(t *testing.T) {
	m := hookModel{
		ObjectKind: "X",
		ObjectAttributes: attributesModel{
			Id:     123,
			Status: "success",
		},
		User: userModel{},
		Project: projectModel{
			Name:   "DSADSADAS",
			WebUrl: "HHHHHHHHHHH",
		},
		Commit: commitModel{
			Message: "DSAVJVJCJCJCJCJCJCJ",
			Author: authorModel{
				Name:  "AAA",
				Email: "BBB",
			},
		},
		Builds: nil,
	}
	text := getText(m)
	fmt.Println(text)
	m.ObjectAttributes.Status = "pending"
	text = getText(m)
	fmt.Println(text)
	m.ObjectAttributes.Status = "running"
	text = getText(m)
	fmt.Println(text)
	m.ObjectAttributes.Status = "failed"
	text = getText(m)
	fmt.Println(text)
	m.ObjectAttributes.Status = "dddd"
	text = getText(m)
	fmt.Println(text)
}
