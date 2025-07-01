package command

import "testing"

func TestSet(t *testing.T) {
	tok := NewTokenizer("SET x \"1\"")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Type != Identifier {
		t.Fatalf("expected token of type Identifier, got %s", token.Type)
	}

	if token.Value != "SET" {
		t.Fatalf("expected token with value \"SET\", got \"%s\"", token.Value)
	}

	token, err = tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Type != Identifier {
		t.Fatalf("expected token of type Identifier, got %s", token.Type)
	}

	if token.Value != "x" {
		t.Fatalf("expected token with value \"x\", got \"%s\"", token.Value)
	}

	token, err = tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Type != String {
		t.Fatalf("expected token of type String, got %s", token.Type)
	}

	if token.Value != "1" {
		t.Fatalf("expected token with value \"1\", got \"%s\"", token.Value)
	}
}

func TestStringScaping(t *testing.T) {
	tok := NewTokenizer("\"{\\\"key\\\": \\\"value\\\\\\\"}\"")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Type != String {
		t.Fatalf("expected token of type Identifier, got %s", token.Type)
	}

	if token.Value != "{\"key\": \"value\\\"}" {
		t.Fatalf("expected token with value \"{\"key\": \"value\\\"}\", got \"%s\"", token.Value)
	}
}

func TestGet(t *testing.T) {
	tok := NewTokenizer("GET x")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Type != Identifier {
		t.Fatalf("expected token of type Identifier, got %s", token.Type)
	}

	if token.Value != "GET" {
		t.Fatalf("expected token with value \"GET\", got \"%s\"", token.Value)
	}

	token, err = tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Type != Identifier {
		t.Fatalf("expected token of type Identifier, got %s", token.Type)
	}

	if token.Value != "x" {
		t.Fatalf("expected token with value \"x\", got \"%s\"", token.Value)
	}
}

func TestEOF(t *testing.T) {
	tok := NewTokenizer("GET x")
	_, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	_, err = tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}

	_, err = tok.NextToken()

	if err == nil {
		t.Fatal("expected error, but tokenizer continued")
	}

	if err.Error() != "EOF" {
		t.Fatalf("expected \"EOF\" error, but got \"%s\"", err.Error())
	}
}
