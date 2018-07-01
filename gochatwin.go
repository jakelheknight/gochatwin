package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type msg struct {
	text        string
	sender      string
	reciever    string // Optional sent only to that single reciever when set.
	channelName string
	timeStamp   time.Time
}

func (message msg) format() string {
	return fmt.Sprintf("\u001b[32m%s\u001b[0m: \u001b[34m(%s)\u001b[0m %s\n\r", message.sender, message.channelName, message.text)
}

type channel struct {
	name       string
	msgStream  chan msg
	subscribed map[string]user
}

func (channel channel) systemMsg(message string) msg {
	return msg{
		text:        message,
		sender:      "System",
		channelName: channel.name,
		timeStamp:   time.Now(),
	}
}

type user struct {
	name    string
	out     chan msg
	focused channel
}

// A few maps of objects registered with the server. When you set a chatManager up it must have a GENERAL channel added to send system level messages.
type chatManager struct {
	users       map[string]user
	channelList map[string]channel
}

func (chatManager *chatManager) makeChannel(channelName string) {
	if _, ok := chatManager.channelList[channelName]; !ok {
		chatManager.channelList[channelName] = channel{
			name:       channelName,
			msgStream:  make(chan msg, 5),
			subscribed: make(map[string]user, 0),
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
	if _, ok := chatManager.channelList[channelName].subscribed[userName]; !ok {
		chatManager.channelList[channelName].subscribed[userName] = chatManager.users[userName]
	}
	// reassigning this user gets aroudn an issue with grabbing an non addressable item inside a map. There are two solutions coppying the item or using pointers.
	chatManager.users[userName] = joinUser

	chatManager.channelList[channelName].msgStream <- chatManager.channelList[channelName].systemMsg(fmt.Sprintf("%s has joint the channel. Say hello.", userName))
}

func (chatManager *chatManager) unJoinChannel(userName string, channelName string) {
	delete(chatManager.channelList[channelName].subscribed, userName)
	if chatManager.users[userName].focused.name == channelName {
		var userObj = chatManager.users[userName]
		userObj.focused = chatManager.channelList["GENERAL"]
		chatManager.users[userName] = userObj
	}
}

func (chatManager *chatManager) handleInput(input string, userName string, channelName string) msg {
	commandArr := strings.Fields(input)
	switch {
	case commandArr[0] == "/help":
		return msg{
			text:        "You can join a channel /join <channel>, unjoin a channel /unjoin <channel> or wisper to any user /w <username>",
			sender:      "System",
			channelName: "GENERAL",
			timeStamp:   time.Now(),
		}
	case commandArr[0] == "/w":
		return msg{
			text:        strings.Join(commandArr[2:], " "),
			sender:      userName,
			reciever:    commandArr[1],
			channelName: "WHISPER",
			timeStamp:   time.Now(),
		}
	case commandArr[0] == "/join":
		chatManager.joinChannel(userName, commandArr[1])
		return msg{
			text:        "You successfully joined a " + commandArr[1],
			sender:      "SYSTEM",
			reciever:    userName,
			channelName: "WHISPER",
			timeStamp:   time.Now(),
		}
	case commandArr[0] == "/unjoin":
		chatManager.unJoinChannel(userName, commandArr[1])
		return msg{
			text:        "You successfully unjoined the channel " + commandArr[1],
			sender:      "SYSTEM",
			reciever:    userName,
			channelName: "WHISPER",
			timeStamp:   time.Now(),
		}
	default:
		return msg{
			text:        input,
			sender:      userName,
			channelName: channelName,
			timeStamp:   time.Now(),
		}
	}
}

func (chatManager *chatManager) run() {
	for {
		for _, channel := range chatManager.channelList {
			message := <-channel.msgStream
			for _, user := range channel.subscribed {
				user.out <- message
			}
		}
	}
}

func handleUserConnection(chatManager *chatManager, conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	var userName string
	io.WriteString(conn, chatManager.channelList["GENERAL"].systemMsg("Welcome to GoChatWin an awesome chat server pleas chose a UserName: ").format())

	// Loop so that users will be unique and not overwritte eachother.
	// Could also postpend a random string of chars at the end but I like this better even if you could use invisable chars this way you could add wisper if you wanted.
	for {
		scanner.Scan()
		userName = scanner.Text()
		if _, ok := chatManager.users[userName]; !ok {

			chatManager.users[userName] = user{
				name:    userName,
				focused: chatManager.channelList["GENERAL"],
			}

			io.WriteString(conn, chatManager.channelList["GENERAL"].systemMsg("Thanks for joining us. Type /help for a list of commands. ").format())

			break
		}
		io.WriteString(conn, chatManager.channelList["GENERAL"].systemMsg("Sorry that user name is taken Please choose another one:").format())
	}

	chatManager.joinChannel(userName, "GENERAL")

	go func() {
		for scanner.Scan() {
			input := scanner.Text()
			chatManager.users[userName].focused.msgStream <- chatManager.handleInput(input, userName, chatManager.users[userName].focused.name)
		}
	}()

	for message := range chatManager.users[userName].out {
		io.WriteString(conn, message.format())
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
		name:       "GENERAL",
		msgStream:  make(chan msg, 5),
		subscribed: make(map[string]user, 0),
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleUserConnection(chatManager, conn)
	}
}
