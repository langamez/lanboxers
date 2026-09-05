package terminal

import (
	"bufio"
	"fmt"
	"github.com/langamez/lanboxers/domain"
	"os"
	"strings"
)

func GetConnectionInfo() (isHost bool, ip string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Host or Join? (h/j): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "h", "host":
			return true, ""

		case "j", "join":
			for {
				fmt.Print("Host IP: ")

				ip, _ = reader.ReadString('\n')
				ip = strings.TrimSpace(ip)

				if ip != "" {
					return false, ip
				}

				fmt.Println("IP cannot be empty.")
			}

		default:
			fmt.Println("Please enter 'h' or 'j'.")
		}
	}
}

func ReadActions(actions chan<- domain.Act) {
	buf := make([]byte, 3)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		act, ok := actionFromKey(string(buf[:n]))
		if !ok {
			continue
		}

		actions <- act
	}
}

func actionFromKey(key string) (domain.Act, bool) {
	switch key {
	case "\033[A":
		return -domain.DoMoveDown, true
	case "\033[B":
		return domain.DoMoveDown, true
	case "\033[C":
		return -domain.DoMoveLeft, true
	case "\033[D":
		return domain.DoMoveLeft, true

	case "n":
		return -domain.DoPunch, true
	case "m":
		return domain.DoPunch, true

	case "q", "Q":
		return domain.Quit, true
	}

	return 0, false
}
