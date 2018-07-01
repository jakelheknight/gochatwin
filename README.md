# GoChatWin

## A basic chat app for the TORBIT Go programming challenge.

The goal of this project is to build a simple chat server in Go.
Multiple clients should be able to connect via telnet and send
messages to the server. When a message is sent to the server it
should be relayed to all the connected clients (including a
timestamp and the name of the client that sent the message). All
messages should also be logged to a local log file. Basic
configuration settings like listening port, ip, and log file location
should be read from a local config file.

## Making it work
``` bash
# First go grab it from get hub or unpack the zipped file in your <GOPATH>/src folder.
go get github.com/jakelhekngiht/gochatwin # I prefer go get.
# Use go install to make the executable.
go install gochatwin
#next add an .env file in the folder you will be running from.
touch .env
# This .env file must contain the following
# PORT=<The desired port number>
# LOG_FILE_LOCATIN=<The folder you want logs printed to>

# Finally just run the built gochat win file.
# To run it in windows):
gochatwin.exe
```

Once you have it up and running jump to another terminal and telnet into localhost PORT#
Your up and running.

## Features
I chose to make it do a little more than just a basic chat. I chose to tackle wispers and multiple chat channels. And pretty it up a bit. 

* You have the ability to wisper to other users in your focused channel.
    * It will send from the current channel from you.
    * If the user is not in the channel or doesnt exist it will do nothing.
* You can join or unjoin any channel you like
    * If you join a new channel one will be created and you will subscibe to and focus that channel.
    * You can subscribe to as many channels as you want.
    * If the channel already exists it will just subscribbe and focus that channel.
    * If you are already subscribed to the channel you will simply focus the channel.
* You can unjoin any channel
    * If you are subscribed to the channel you will be removed from the map.
    * If you are fouced to the channel your focuse will be pointed at the GENERAL channel.
* When user leaves and comes back he is still subscribed to the previous channels but is focuse on the general.

## Unaddressed Issues
Since this was a limited exercise I didn't go thrugh and build more than the basics above and left off some steps I would do for production code. This is not exhastive.

* No sanitization of anything. For a more production moddel you should never trust the user to behave.
* Users input is overwritten as message comes in. If I were to fix this I would make a user client to handle the inputs and such.
* If you chat too soon while the log file is being opened you can end up missing the first message or two. This could be fixed by opening the log file in the main function and passing it rather than opening it before the chatManatger loop.
* Ther is no channel clean up. As time goes on you will end up with more and more channels untill you run out of memory. This issue is small but could be solved with an empty check on unjoin that would delete the chat channel when empty.
* I would write a few more test if this were going to be production level.

## Creadit
* I used this awesome toutorial as a startring point. 
    * https://www.youtube.com/watch?v=cNxfgXrHeAg
    * https://github.com/golang-book/bootcamp-examples/blob/master/week2/day2/chat-server/main.go
* I also used godotenv
    * https://github.com/joho/godotenv