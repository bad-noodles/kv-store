package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bad-noodles/kv-store/pkg/client"
	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
)

func humanFriendly(v typesystem.Type) string {
	switch value := v.(type) {
	case typesystem.Status:
		return value.Value().(string)
	case typesystem.StringValue:
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(value.Value().(string), "\"", "\\\""))
	case typesystem.ArrayValue:
		var b strings.Builder
		items := value.Value().([]typesystem.Type)
		last := len(items) - 1

		b.WriteString("[ ")

		for i, v := range items {
			b.WriteString(humanFriendly(v))
			if i != last {
				b.WriteString(", ")
			}
		}

		b.WriteString(" ]")
		return b.String()
	default:
		return fmt.Sprint(value)
	}
}

func main() {
	client := client.NewClient()
	client.Connect("localhost:1337")
	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			break
		}
		input := strings.Trim(line, " \n")

		if input == "exit" {
			break
		}

		err = client.Execute(input)
		if err != nil {
			panic(err)
		}

		resp, err := client.Read()
		if err != nil {
			panic(err)
		}

		fmt.Println(humanFriendly(resp))

	}
}
