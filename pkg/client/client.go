package client

import (
	"fmt"
	"net"

	"github.com/bad-noodles/kv-store/pkg/command"
	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
)

type Client struct {
	conn          net.Conn
	typeParser    *typesystem.Parser
	commandParser *command.Parser
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	c.conn = conn
	c.typeParser = typesystem.NewParser(c.conn)
	c.commandParser = command.NewParser()

	return err
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Execute(input string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	cmd, err := c.commandParser.Parse(input)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(c.conn, cmd.String())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Read() (typesystem.Type, error) {
	c.typeParser.Next()

	return c.typeParser.Data(), c.typeParser.Error()
}
