# Distributed Social Network

Distributed Social Network (yet to be named) is a MakeBU project.

## Authors:

* Emmanuel Amponsah
* Hayato Nakamura
* Joseph Lai
* Ken Garber
* Shoki Ko
* William Pine

## Flowchart Overview
* User starts the program (their server), and it either creates its own network or joins a specified network.
* Server goes into an infinite loop to support two services:
  * 1: network communications, and
  * 2: user interface
* User interface allows the user to:
  * 1: view all chatrooms, view and send messages in chatrooms, and
  * 2: view and set their identity.
* Database layer supports the server side in order to store data persistently and recall data

## Supporting Tools
* Python server configuration and load visualization
* C++ database compression tool

## Components
* Server
  * Network Communication
  * Local User Interface
* User Interface
  * Local server
  * Front end (in browser)
* Database layer

