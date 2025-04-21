package main

import (
	"fmt"
	"io"
	"os"
)

const chunkSize = 50 * 1024 // 50 kb

func main()  {
	filePath := "testfile.mp4"

	file, err := os.Open(filePath); 
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

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

		chunkFileName := fmt.Sprintf("chunk_%d", chunkCount)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			panic(err)
		}

		_, err = chunkFile.Write(buffer[:byteRead])
		if err != nil {
			panic(err)
		}

		chunkFile.Close()
		chunkCount++
	}

	fmt.Printf("✅ File split into %d chunks.\n", chunkCount)

	restruct()
}