package server

import (
	"fmt"
	"log"
	"net"

	"github.com/bad-noodles/kv-store/pkg/store"
	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
)

func handleConnection(conn net.Conn, st *store.Store) {
	defer conn.Close()

	parser := typesystem.NewParser(conn)
	fmt.Println(1)

	for parser.Next() {
		fmt.Println(2)
		if parser.Error() != nil {
			_, err := fmt.Fprint(conn, typesystem.NewStatus(false, parser.Error().Error()))
			if err != nil {
				log.Println(err)
			}
		}

		data := parser.Data()
		fmt.Println("query")

		switch cmd := data.(type) {
		case typesystem.ArrayValue:
			resp := st.ExecuteParsedCommand(cmd)

			_, err := fmt.Fprint(conn, resp.String())
			if err != nil {
				log.Println(err)
			}
		default:
			_, err := fmt.Fprint(conn, typesystem.NewStatus(false, "Not a command"))
			if err != nil {
				log.Println(err)
			}
		}
	}

	fmt.Println("Client disconnected")
}

func Start(port int) {
	walPath := "./wal"
	st := store.NewStore(walPath)
	st.Restore(walPath)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on port %d\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Default().Print(err)
		}

		fmt.Println("Client connected")

		go handleConnection(conn, st)
	}
}
