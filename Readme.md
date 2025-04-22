# 📦 MediaShare

> This project is a secure file transfer system built in Go, designed to send and receive files over a local network with AES-GCM encryption. It mimics tools like Xender or AirDrop, but runs from the command line and offers encryption-by-default to protect data during transmission.  


## ✨ Features

- ✅ File transfer over TCP
- ⚡ CLI-based interface for sending and receiving
- 🔒 AES-256 GCM encryption with unique nonce per chunk

## 🚀 Demo

[Live Demo]("https://www.youtube.com/watch?v=LwRR86JO5Os") 

## 🔧 Usage
Run Receiver:
```
go run main.go --mode=receive
```
Run Sender:
```
go run main.go --mode=send --filePath=/path/to/file --target=127.0.0.1
```
Replace 127.0.0.1 with the receiver's IP