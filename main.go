package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"flag"
	"net"
	"encoding/binary"
	"path/filepath"
)

const chunkSize = 50 * 1024 // 50 kb
const port = 9002
var key = []byte("thisis32byteslongpassphrase!!__i")
const nonceSize = 12

func send(receiverIP, filePath string) {
	target := net.JoinHostPort(receiverIP, fmt.Sprint(port))
	conn, err := net.Dial("tcp", target)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	file, err := os.Open(filePath); 
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	filename := filepath.Base(filePath)
	filenameBytes := []byte(filename)
	filenameLen := uint32(len(filenameBytes))

	filenameLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(filenameLenBuf, filenameLen)

	// send file metadata(name, chunksize)
	conn.Write(filenameLenBuf)
	conn.Write(filenameBytes)

	buffer := make([]byte, chunkSize)
	chunkCount := 0

	for {
		byteRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if byteRead == 0 {
			break // EOF
		}

		chunkData := buffer[:byteRead]
		chunkSizeBytes := make([]byte, 4) // 
		
		// encrypt
		encryptedData, err := encryptAESGCM(chunkData, key)
		if err != nil {
			panic(err)
		}
		// 🔹 Then send the encrypted chunk
		binary.BigEndian.PutUint32(chunkSizeBytes, uint32(len(encryptedData)))
		conn.Write(chunkSizeBytes)
		
		_, err = conn.Write(encryptedData)
		if err != nil {
			fmt.Println("❌ Failed to send encrypted chunk:", err)
			return
		}

		chunkCount++
	}

	fmt.Printf("✅ Encrypted File split into %d chunks.\n", chunkCount)

}

func receive() {
	// extract filename from conn
	// extract size, chunkSize
	// extract encrypted data
	// resemble the file

	listener, err := net.Listen("tcp", ":9002")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("📥 Listening for incoming file...")

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// extract filename from conn
	filenameLenBuf := make([]byte, 4)
	_, err = io.ReadFull(conn, filenameLenBuf)
	if err != nil {
		panic("❌ Failed to read filename length: " + err.Error())
	}
	filenameLen := binary.BigEndian.Uint32(filenameLenBuf)

	// Read the filename itself
	filenameBuf := make([]byte, filenameLen)
	_, err = io.ReadFull(conn, filenameBuf)
	if err != nil {
		panic("❌ Failed to read filename: " + err.Error())
	}
	filename := string(filenameBuf)
	fmt.Println("📁 Receiving file:", filename)

	outputFile, err := os.Create("received_" + filename)
	if err != nil {
		panic("❌ Could not create output file: " + err.Error())
	}
	defer outputFile.Close()

	fmt.Println("📥 Receiving from", conn.RemoteAddr())

	for {
		sizeBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, sizeBuf)
		if err == io.EOF{
			break
		}

		if err != nil {
			fmt.Println("❌ Error reading chunk size:", err)
			break
		}

		chunkSize := binary.BigEndian.Uint32(sizeBuf)
		encData := make([]byte, chunkSize)
		_, err = io.ReadFull(conn, encData)

		if err != nil {
			fmt.Println("❌ Error reading encrypted chunk:", err)
			break
		}

		plainData, err := decryptAESGCM(encData, key)
		if err != nil {
			fmt.Println("❌ Failed to decrypt:", err)
			break
		}

		_, err = outputFile.Write(plainData)
		if err != nil {
			fmt.Println("❌ Failed to write to file:", err)
			break
		}
	}
}

func encryptAESGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	cipherText := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, cipherText...), nil
}

func decryptAESGCM(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce := ciphertext[:nonceSize]
	ciphertextData := ciphertext[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, nonce, ciphertextData, nil)
}

func main()  {
	// mode, filepath, targetIP

	var mode = flag.String("mode", "", "To indicate the operation mode (send or receive)")
	var filePath = flag.String("filePath", "", "Path to the file you want to send")
	var receiverIP = flag.String("target", "", "IP address of the receiver")

	flag.Parse()

	if *mode == "send" && *filePath != "" && *receiverIP != "" {
		send(*receiverIP, *filePath)
	}else if *mode == "receive" {
		receive()
	}else {
		fmt.Printf("You have likely not used the correct command, use --help to see available options")
	}
}
