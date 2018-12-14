# Chat on Go

This is just a fun project as I learn go. I'll be implementing chat functionality through various mediums e.g. file, command line, web & chatbot

## File chat
First version of the app (no tests yet).
To try it out switch to the file directory then run `go run main.go`
Respond with the list of users taking part in the chat e.g. `user1 user2 user3`. 

There's really no limit on the users at the moment.
Open any of the txt files created e.g. user1.txt and type your message then save the file. The message should appear on the other files for the users created.

Due to constraints in reloading files in an editor, only one user can send a message at a time. Obviously this would work better on a normal UI, but since we are working within the constraints of a file  it's not too horrible.

## TBD
- possibly track file changes so as to keep the chat consistent
- command line chat: still trying to think how to implement that
- web: react app with go backend
- chatbot: similar to web, but will have an API that can be extended to use e.g. on FB messenger, Telegram etc
Sounds ironical to build a chat app on top of another chat app right? Oh well! It's just for fun.