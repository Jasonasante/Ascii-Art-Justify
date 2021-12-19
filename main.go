package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// this function retrieves the width of the terminal
func getWidth() uint {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col)
}

//to lower case the assigned alignment in case user uses a capital letter
func ToLower(s string) string {
	lowstr := []rune(s)
	for i, char := range lowstr {
		if char >= 65 && char <= 90 {
			lowstr[i] = lowstr[i] + 32
		}
	}
	return string(lowstr)
}

//print the incoming string as a GUI in the preferred alignment stated by the user.
func main() {
	var emptyString string
	var inputString []string
	if len(os.Args) == 4 {
		inputString = strings.Split(os.Args[1], "\\n")
	} else {
		fmt.Println("Usage: go run . [STRING] [BANNER] [OPTION]")
		fmt.Println("EX: go run . something standard --align=right")
		os.Exit(0)
	}
	Content, _ := os.ReadFile(os.Args[2] + ".txt")
	asciiSlice2 := make([][]string, 95)
	// this stores the ascii-bubbles in order of the
	// there are 95 ascii characters and this lets us index the dimension holding each bubble
	for i := 0; i < len(asciiSlice2); i++ {
		asciiSlice2[i] = make([]string, 9)
	}
	// this makes the asciiSlice2[i] have a length of 8
	var bubbleCount int
	count := 0
	for i := 1; i < len(Content); i++ {
		if Content[i] == '\n' && bubbleCount <= 94 {
			asciiSlice2[bubbleCount][count] = emptyString
			// so bubbleCount is the index and count is the row
			// so asciiSlice2[1][0] is the 1st row of the exclamation mark
			emptyString = ""
			count++
		}
		if count == 9 {
			count = 0
			bubbleCount++
		} else {
			if Content[i] != '\n' && Content[i] != '\r' {
				emptyString += string(Content[i])
				// as count != 8 and Contet[i]!= '\n', it will append the contents of each row.
				// Once it reaches the '\n' at the end of the row, the first if statement is activated.
			}
		}
	}

	var alignFlag []string
	estrCount := 0
	var tempOutput [][]string

	if strings.HasPrefix(os.Args[3], "--align=") {
		alignFlag = strings.Split(os.Args[3], "--align=")
	} else {
		fmt.Println("Usage: go run . [STRING] [BANNER] [OPTION]")
		fmt.Println("EX: go run . something standard --align=right")
		os.Exit(0)
	}

	alignFlag[1] = ToLower(alignFlag[1])

	length := int(getWidth())
	charlength := 0
	startingPoint := 0
	for j, str := range inputString {
		for _, aRune := range str {
			tempOutput = append(tempOutput, asciiSlice2[aRune-rune(32)])
			// due to the loop it will append the bubble eqivalent of the every letter inside inputString
		}
		// the loop below is to get the length of the first line of every bubble character
		for h := 0; h < len(tempOutput); h++ {
			charlength += len(tempOutput[h][0])
		}
		for i := range tempOutput[0] {
			for _, char := range tempOutput {
				if alignFlag[1] == "center" && estrCount == 0 {
					startingPoint = (length / 2) - (charlength / 2)
					for l := 0; l < startingPoint; l++ {
						fmt.Printf(" ")
					}
					fmt.Print(char[i])
				} else if alignFlag[1] == "left" {
					fmt.Print(char[i])
				} else if alignFlag[1] == "right" && estrCount == 0 {
					startingPoint = (length - charlength)
					for l := 0; l < startingPoint; l++ {
						fmt.Printf(" ")
					}
					fmt.Print(char[i])
				} else {
					fmt.Print(char[i])
				}
				//now figure how how to replace %110v with the
				estrCount++
				if estrCount == len(inputString[j]) {
					estrCount = 0
				}
				// // this prints each line of each bubble letter (which is each slice of string)
			}
			fmt.Println()
		}
		//fmt.Println(charlength)
		tempOutput = nil
		charlength = 0
		// once the word has been printed, we want to reset tempOutput to nil, ready to be filled
		// by the next string element in inputString.
	}
}
