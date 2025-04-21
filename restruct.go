package main

import (
	"fmt"
	"io"
	"os"
)

func restruct(){
	// Open the file
	outputFile, err := os.Create("retestfile.mp4")

	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	for i := 0; ; i++ {
		chunkName := fmt.Sprintf("chunk_%d", i)
		chunk, err := os.Open(chunkName)

		if os.IsNotExist(err) {
			break
		}

		if err != nil {
			panic(err)
		}

		_, err = io.Copy(outputFile, chunk)
		if err != nil {
			panic(err)
		}
		chunk.Close()
		fmt.Println("Added:", chunkName)
	}
	fmt.Println("✅ Reassembled file saved as:", outputFile)
}