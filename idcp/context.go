package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Context struct {
	FunctionName     string
	Changed          bool
	Comment          string
	ErrString        string
	AdditionalValues map[string]string
}

func NewContext(functionName string) *Context {
	return &Context{
		FunctionName:     functionName,
		Changed:          false,
		Comment:          "",
		ErrString:        "",
		AdditionalValues: make(map[string]string),
	}
}

func (c *Context) AddComment(comment string) {
	if len(c.Comment) > 0 {
		c.Comment += " | "
	}
	c.Comment += comment
}

func (c *Context) AddValue(key, value string) {
	c.AdditionalValues[key] = value
}

func (c *Context) GetValue(key string) string {
	value, _ := c.AdditionalValues[key]
	return value
}

func (c *Context) SetError(comment string, err error) {
	c.ErrString = comment + err.Error()
	c.Finish()
}

func (c *Context) Begin() {
	EmitTimestamp("start")
	EmitValue("function", c.FunctionName)
}

func (c *Context) Finish() {
	for key, value := range c.AdditionalValues {
		EmitValue(key, value)
	}

	EmitChanged(c.Changed)
	EmitValue("comment", c.Comment)
	EmitValue("error", c.ErrString)
	EmitTimestamp("finish")

	if len(c.ErrString) != 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

func EmitValue(key, value string) {
	_, err := fmt.Printf("%-20s%s\n", key, value)
	if err != nil {
		log.Fatal("Error emitting", key, "value:", err.Error())
	}
	return
}

func EmitTimestamp(key string) {
	timestamp := time.Now().Format(time.RFC3339)
	EmitValue(key, timestamp)
	return
}

func EmitChanged(changed bool) {
	changedStr := "false"
	if changed {
		changedStr = "true"
	}

	EmitValue("changed", changedStr)
}
