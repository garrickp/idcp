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
	c.EmitValue("start", "")
	c.EmitValue("function", c.FunctionName)
}

func (c *Context) Finish() {
	for key, value := range c.AdditionalValues {
		c.EmitValue(key, value)
	}

	c.EmitChanged(c.Changed)
	c.EmitValue("comment", c.Comment)
	c.EmitValue("error", c.ErrString)
	c.EmitValue("finish", "")

	if len(c.ErrString) != 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

func (c *Context) EmitValue(key, value string) {
	timestamp := time.Now().Format(time.RFC3339)
	_, err := fmt.Printf("%s\t%s\t%-12s\t%s\n", timestamp, c.FunctionName, key, value)
	if err != nil {
		log.Fatal("Error emitting", key, "value:", err.Error())
	}
	return
}

func (c *Context) EmitChanged(changed bool) {
	changedStr := "false"
	if changed {
		changedStr = "true"
	}

	c.EmitValue("changed", changedStr)
}
