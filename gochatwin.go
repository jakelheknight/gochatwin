package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type msg struct {
	text      string
	sender    string
	reciever  string // Optional sent only to that single reciever when set.
	timeStamp time.Time
}

func (message msg) format(channel string) string {
	return fmt.Sprintf("\u001b[32m%s\u001b[0m: \u001b[30m(%s)\u001b[0m %s\n", message.sender, channel, message.text)
}

type channel struct {
	name      string
	msgStream chan msg
}

func (channel channel) systemMsg(message string) msg {
	return msg{
		text:      message,
		sender:    "System",
		timeStamp: time.Now(),
	}
}

type user struct {
	name       string
	focused    channel
	subscribed map[string]channel
}

// A few maps of objects registered with the server. When you set a chatManager up it must have a GENERAL channel added to send system level messages.
type chatManager struct {
	users       map[string]user
	channelList map[string]channel
}

func (chatManager *chatManager) makeChannel(channelName string) {
	if _, ok := chatManager.channelList[channelName]; !ok {
		chatManager.channelList[channelName] = channel{
			name:      channelName,
			msgStream: make(chan msg, 5),
		}
		chatManager.channelList["GENERAL"].msgStream <- chatManager.channelList[channelName].systemMsg(fmt.Sprintf("New Channel: %s is ready for use.", channelName))
	} else {
		chatManager.channelList["GENERAL"].msgStream <- chatManager.channelList[channelName].systemMsg(fmt.Sprintf("Channel: %s already exists.", channelName))
	}
}

func (chatManager *chatManager) joinChannel(userName string, channelName string) {
	if _, ok := chatManager.channelList[channelName]; !ok {
		chatManager.makeChannel(channelName)
	}
	joinUser := chatManager.users[userName]
	joinUser.focused = chatManager.channelList[channelName]
	if _, ok := chatManager.users[userName].subscribed[channelName]; !ok {
		chatManager.users[userName].subscribed[channelName] = chatManager.channelList[channelName]
	}
	// reassigning this user gets aroudn an issue with grabbing an non addressable item inside a map. There are two solutions coppying the item or using pointers.
	chatManager.users[userName] = joinUser

	chatManager.channelList[channelName].msgStream <- chatManager.channelList[channelName].systemMsg(fmt.Sprintf("%s has joint the channel. Say hello.", userName))
}

var general = channel{
	name:      "General",
	msgStream: make(chan msg),
}

func handleUserConnection(chatManager *chatManager, conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	var userName string
	io.WriteString(conn, chatManager.channelList["GENERAL"].systemMsg("Welcome to GoChatWin an awesome chat server pleas chose a UserName: ").format("GENERAL"))

	// Loop so that users will be unique and not overwritte eachother.
	// Could also postpend a random string of chars at the end but I like this better even if you could use invisable chars this way you could add wisper if you wanted.
	for {
		scanner.Scan()
		userName = scanner.Text()
		if _, ok := chatManager.users[userName]; !ok {

			chatManager.users[userName] = user{
				name:       userName,
				focused:    general,
				subscribed: make(map[string]channel, 0),
			}

			io.WriteString(conn, chatManager.channelList["GENERAL"].systemMsg("Thanks for joining us. Type /help for a list of commands. ").format("GENERAL"))

			break
		}
		io.WriteString(conn, chatManager.channelList["GENERAL"].systemMsg("Sorry that user name is taken Please choose another one:").format("GENERAL"))
	}

	chatManager.joinChannel(userName, "GENERAL")

	go func() {
		for scanner.Scan() {
			input := scanner.Text()

			chatManager.users[userName].focused.msgStream <- msg{
				text:      input,
				sender:    chatManager.users[userName].name,
				timeStamp: time.Now(),
			}
		}
	}()

	for _, channel := range chatManager.users[userName].subscribed {
		for message := range channel.msgStream {
			io.WriteString(conn, message.format(channel.name))
		}
	}
}

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	chatManager := &chatManager{
		users:       make(map[string]user, 0),
		channelList: make(map[string]channel, 0),
	}

	chatManager.channelList["GENERAL"] = channel{
		name:      "GENERAL",
		msgStream: make(chan msg, 5),
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleUserConnection(chatManager, conn)
	}
}
