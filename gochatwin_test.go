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

// Need new object for unit testing specific functions.
var testChatManager = chatManager{
	users:       make(map[string]*user, 0),
	channelList: make(map[string]*channel, 0),
	msgStream:   make(chan msg, 5),
}

var testUers = user{
	name:    "Test",
	focused: "GENERAL",
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
	testChatManager.makeChannel("Test")
	if _, ok := testChatManager.channelList["Test"]; !ok {
		t.Logf("Failed to creat cannel when make chanel was called.")
	}
}

func TestChatManagerJoin(t *testing.T) {
	testChatManager.joinChannel("Test", "NewChannel")

}
