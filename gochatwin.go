package main

import (
	"bufio"
	"fmt"
	"io"
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
	users     map[string]user // Key is to match user.name
}

type user struct {
	name    string
	focused channel
}

// A few maps of objects registered with the server. When you set a server up it must have a GENERAL channel added to send system level messages.
type server struct {
	users       map[string]user
	channelList map[string]channel
}

func (server *server) makeChannel(channelName string) {
	if _, ok := server.channelList[channelName]; !ok {
		server.channelList[channelName] = channel{
			name:      channelName,
			msgStream: make(chan msg, 0),
			users:     make(map[string]user, 0),
		}
		server.channelList["GENERAL"].msgStream <- msg{
			text:      fmt.Sprintf("New Channel: %s is ready for use.", channelName),
			sender:    "System",
			timeStamp: time.Now(),
		}
	} else {
		server.channelList["GENERAL"].msgStream <- msg{
			text:      fmt.Sprintf("Channel: %s already exists.", channelName),
			sender:    "System",
			timeStamp: time.Now(),
		}
	}
}

func (server *server) joinChannel(userName string, channelName string) {
	if _, ok := server.channelList[channelName]; ok {
		joinUser := server.users[userName]
		joinUser.focused = server.channelList[channelName]
		if _, ok := server.users[userName]; !ok {
			server.channelList[channelName].users[userName] = joinUser
		}
		// reassigning this user gets aroudn an issue with grabbing an non addressable item inside a map. There are two solutions coppying the item or using pointers.
		server.users[userName] = joinUser
	}

}

var general = channel{
	name:      "General",
	msgStream: make(chan msg),
	users:     make(map[string]user, 0),
}

func handleUserConnection(server *server, conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	var currentUser user
	var welcomeMsg = msg{
		text:      "Welcome to GoChatWin an awesome chat server pleas chose a UserName: ",
		sender:    "System",
		timeStamp: time.Now(),
	}
	var userSucessMsg = msg{
		text:      "Thanks for joining us. Type /help for a list of commands. ",
		sender:    "System",
		timeStamp: time.Now(),
	}
	var userFailMsg = msg{
		text:      "Sorry that user name is taken Please choose another one:",
		sender:    "System",
		timeStamp: time.Now(),
	}
	io.WriteString(conn, welcomeMsg.format("GENERAL"))

	// Loop so that users will be unique and not overwritte eachother.
	// Could also pull a bungy and postpend a random string of chars at the end but I like this better even if you could use invisable chars.
	for {
		scanner.Scan()
		userName := scanner.Text()
		if _, ok := server.users[userName]; ok {
			currentUser = user{
				name:    userName,
				focused: general,
			}
			server.users[userName] = currentUser

			io.WriteString(conn, userSucessMsg.format("GENERAL"))

			break
		}
		io.WriteString(conn, userFailMsg.format("GENERAL"))
	}

	server.joinChannel(currentUser.name, "GENERAL")

}

func main() {

}
