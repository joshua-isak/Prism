# Prism

An End to End encrypted chat room written in go with a custom communication protol (Prism Protocol) that allows for cross platform client development.

## Getting Started

Simply download the binary relevant to your operating system and run it!
Enter the server IP address and a username of your choice. You'll also need a 32-Byte key that other people in the server have in order to see their messages. Only clients with the same 32-Byte key can decrypt each other's messages.

## Prism Protocol (PP)

Prism Protocol is a low level application protocol built on top of TCP that governs communications between Prism servers and clients.

PP uses the concept of "packet types" to divide communications into disrete functions. The currently implemented "packet pypes" are as follows:

- Initial
- Welcome
- General Message

### Initial [1]

After establishing a TCP connection to the server, the first thing a client sends to the server is an "Initial" type packet. This packet contains the packet type (1 for the Initial type), the length of the client's username, the client's username, as well as the client's version. If the client doesn't send this packet within 20 seconds the TCP connection is closed.

Once the Initial packet is sent the server will make sure the client's version is compatible, and that the data sent makes sense. If either of those cases are not true then the client will be disconnected.

NOTE: PP uses little endian byte ordering

| Byte # | Details |
| ------ | ------- |
| 0 | The packet type 1 (uint8) |
| 1 | The length of the username (uint8) |
| 2 to 22 | The username encoded in UTF-8 |
| 23 to 30 | Client version information encoded in UTF-8 |

### Welcome [2]

After sending the Initial packet the client will wait for the server to respond with the Welcome packet. The Welcome packet contains information needed to get the client up to speed with dynamic information. Currently this information is simply the number of users connected to the server as well as their usernames. A maximum of 255 clients can be connected to a single server.

NOTE: Bytes 1 and onward can be empty if there are no users connected to the server

| Byte # | Details |
| ------ | ------- |
| 0 | The packet type 2 (uint8) |
| 1 | Number of currently connected clients (uint8) |
| 2 | The length x of a username (uint8) |
| 3 to x | The username encoded in UTF-8 |
| x + 1 | The length y of a username |
|x + 2 to y | The username encoded in UTF-8 |
| . . . | And so on... |

### GeneralMessage [5]

Once the Welcome packet is received by the client, it can begin sending messages to the server that will be sent to all clients in the room using the GeneralMessage packet type. The server will also send this same packet type back to the client to inform it of messages sent from other users connected to the server.

NOTE: Every client that sends a GeneralMessage receives a reply from the server in the form of a GeneralMessage that is exactly the same as the one it sent out.

NOTE II: Unencryted messages are encoded in UTF-8, but the client should never send an unencrypted message.

NOTE III: Messages are encrypted using 256-bit AES

| Byte # | Details |
| ------ | ------- |
| 0 | The packet type 5 (uint8) |
| 1 | The length of the username (uint8) |
| 2 to 22 | The username encoded in UTF-8 |
| 23 | A boolean, true if the message is encrypted |
| 24 | The length x of the message |
| 25 - x | The message
