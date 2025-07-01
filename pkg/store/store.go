package store

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bad-noodles/kv-store/pkg/command"
	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
)

type Mode int

const (
	Default Mode = iota
	Restore
)

type Store struct {
	data     map[string]typesystem.Type
	wal      *os.File
	mode     Mode
	mutex    sync.RWMutex
	walMutex sync.Mutex
}

func NewStore(walPath string) *Store {
	wal, err := os.OpenFile(walPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	return &Store{
		data: make(map[string]typesystem.Type),
		wal:  wal,
		mode: Default,
	}
}

func (s *Store) ExecuteCommand(input string) typesystem.Type {
	p := command.NewParser()
	cmd, err := p.Parse(input)
	if err != nil {
		return typesystem.NewStatus(false, err.Error())
	}

	return s.ExecuteParsedCommand(cmd)
}

func (s *Store) ExecuteParsedCommand(cmd typesystem.ArrayValue) typesystem.Type {
	args := cmd.Value().([]typesystem.Type)

	switch strings.ToUpper(args[0].Value().(string)) {
	case "SET":
		go s.writeAhead(cmd.String())
		return s.set(args[1:]...)
	case "GET":
		return s.get(args[1])
	default:
		return typesystem.NewStatus(false, "ERROR: Command %v not implemented")
	}
}

func (s *Store) writeAhead(input string) {
	s.walMutex.Lock()
	defer s.walMutex.Unlock()

	if s.mode == Restore {
		return
	}
	fmt.Fprint(s.wal, input)
}

func (s *Store) set(args ...typesystem.Type) typesystem.Type {
	if len(args) < 2 {
		return typesystem.NewStatus(false, "SET requires at least 2 arguments")
	}
	key := args[0]

	switch key.(type) {
	case typesystem.StringValue:
		break
	default:
		return typesystem.NewStatus(false, fmt.Sprintf("Invalid key \"%s\"", key.Value()))
	}

	if s.mode == Default {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}

	s.data[key.Value().(string)] = args[1]

	return typesystem.NewStatus(true, "OK")
}

func (s *Store) get(key typesystem.Type) typesystem.Type {
	switch key.(type) {
	case typesystem.StringValue:
		break
	default:
		return typesystem.NewStatus(false, fmt.Sprintf("Invalid key \"%s\"", key.Value()))
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, found := s.data[key.Value().(string)]

	if !found {
		return typesystem.NewNull()
	}

	return value
}

func (s *Store) Restore(walPath string) {
	f, err := os.Open(walPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s.mutex.Lock()
	s.mode = Restore

	parser := typesystem.NewParser(f)

	for parser.Next() {
		s.ExecuteParsedCommand(parser.Data().(typesystem.ArrayValue))
	}

	if parser.Error() != nil {
		panic(parser.Error())
	}

	s.walMutex.Lock()
	s.mode = Default
	s.walMutex.Unlock()
	s.mutex.Unlock()
}
