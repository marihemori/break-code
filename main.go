package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/yeka/zip"
)

const (
	zipPath          = "./assets/protected.zip"
	pathFilePassword = "./assets/rockyou.txt"
	threadNumber     = 3
	linesPerThread   = 3000
)

func main() {
	openZipFile()
}

func bruteForce(zipPath string, pwdList []string, channel chan<- string) {
	file, err := zip.OpenReader(zipPath)

	if err != nil {
		panic("Error!")
	}

	defer file.Close()

	zipFile := file.File[0]

	for _, value := range pwdList {
		zipFile.SetPassword(string(value))
		_, err := zipFile.Open()

		if err == nil {
			fmt.Print("We found the password!\n")

			zipReader, err := zipFile.Open()

			if err != nil {
				panic("Error!")
			}

			buf, err := ioutil.ReadAll(zipReader)
			if err != nil {
				panic("Error!")
			}

			defer zipReader.Close()

			fmt.Printf("Size of %v: %v byte(s)\n", zipFile.Name, len(buf))

			channel <- string(value)
			break
		}
	}
}

func openZipFile() {
	files, err := zip.OpenReader(zipPath)

	if err != nil {
		panic(fmt.Sprintf("An error occurred: %v", err))
	}

	defer files.Close()

	zipFile := files.File[0]

	if zipFile.IsEncrypted() {
		pwdList := getListOfPasswords(pathFilePassword)
		fmt.Print("The file is proteted by password! \n")

		channel := make(chan string, 1)

		start := time.Now()

		initialLine := 0

		for i := 0; i < threadNumber; i++ {
			finalLine := linesPerThread * (i + 1)

			// outputColor := RandomOutputColor(i)
			// output := color.new(outputColor)
			fmt.Printf("Starting thread %d reading from line %d till line %d\n", i+1, initialLine, finalLine)
			go bruteForce(zipPath, pwdList[initialLine:finalLine], channel)

			initialLine = finalLine + 1
		}

		fmt.Println("--------------------------------------------------------")

		color.Yellow("Cracking the password ...")
		fmt.Println()

		select {
		case password := <-channel:
			fmt.Printf("\nThe password is:\"%v\"\n", password)
			fmt.Printf("\nThe cracking time process took %v \n", time.Since(start))
		case <-time.After(time.Duration(15) * time.Second):
			// fmt.Printf("Timeout after: %d seconds \n", timeout)
			fmt.Printf("Password not found :( \n")
		}
		fmt.Printf("Waiting...")
	}
}

func getListOfPasswords(pathFilePassword string) []string {
	file, error := os.Open(pathFilePassword)

	if error != nil {
		panic(fmt.Sprintf("An error occurred: %v", error))
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var passwords []string

	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
	}

	fileStatus, error := os.Stat(pathFilePassword)

	if error != nil {
		panic(fmt.Sprintf("An error occurred: %v", error))
	}

	fmt.Printf("Status of the file => Name: %v, Size: %v KB \n", fileStatus.Name(), fileStatus.Size()/(1024))
	fmt.Printf("Total number of passwords: %d \n", len(passwords))

	return passwords
}
