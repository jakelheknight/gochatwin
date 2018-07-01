package main

import (
	"strings"
	"testing"
	"time"
)

// Instantiation of objects to be used in tests.
var testMessage1 = msg{
	text:      "Test Message 1",
	sender:    "Test",
	timeStamp: time.Now(),
}

var testUser = user{
	name:    "Test",
	focused: "GENERAL",
}

// Setup channel for testing.

func preChatSetup() *chatManager {
	testChatManager := &chatManager{
		users:       make(map[string]*user, 0),
		channelList: make(map[string]*channel, 0),
		msgStream:   make(chan msg, 5),
	}
	testChatManager.channelList["GENERAL"] = &channel{
		name:       "GENERAL",
		subscribed: make(map[string]string, 0),
	}
	testChatManager.users["Test"] = &testUser
	return testChatManager
}

func TestMsg(t *testing.T) {
	// Since we may want to change the formatting I just need to make sure the actuall message and user are sent.
	switch {
	case !strings.Contains(testMessage1.format(), testMessage1.text):
		t.Logf("You definatly need to have the test in the formatted string.")
	case !strings.Contains(testMessage1.format(), testMessage1.sender):
		t.Logf("You should have the name of the sender in the formatted string.")
	}
}

func TestChatManagerMakeChannel(t *testing.T) {
	testChatManager := preChatSetup()
	testChatManager.makeChannel("Test")
	if _, ok := testChatManager.channelList["Test"]; !ok {
		t.Logf("Failed to creat cannel when make chanel was called.")
	}
}

func TestChatManagerJoin(t *testing.T) {
	testChatManager := preChatSetup()
	testChatManager.joinChannel("Test", "NewChannel")
	if _, ok := testChatManager.channelList["Test"]; !ok {
		t.Logf("Failed to creat cannel when joining a new chanel.")
	}
	if usr, ok := testChatManager.users["Test"]; !ok {
		t.Logf("Failed to put test user on chatmanager.")
	} else if usr.focused != "Test" {
		t.Logf("Failed to subscribe to joined channel.")
	}
}
