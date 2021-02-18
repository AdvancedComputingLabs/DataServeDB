package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cliProcessor() {
	reader := bufio.NewReader(os.Stdin)
	keep_running := true

	for keep_running {

		fmt.Println("Enter Command: ")
		cmd_text, e := reader.ReadString('\n')
		cmd_text = strings.Trim(cmd_text, "\r\n")
		cmd_text_toks := strings.Split(cmd_text, " ")

		if e != nil {
			fmt.Println(e.Error())
			return
		}

		if len(cmd_text_toks) > 0 {
			switch strings.ToUpper(cmd_text_toks[0]) {

			case "EXIT":
				for _, s := range servers {
					s.Shutdown()
				}
				keep_running = false
			}
		}

	}
}
