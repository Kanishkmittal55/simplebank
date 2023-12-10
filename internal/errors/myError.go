package errors

import (
	"errors"
	"fmt"
)

func foo() error {
	return errors.New("FOO had an error")
}

func boo() error {
	err := foo()
	if err != nil {
		return fmt.Errorf("boo had an error : %w", err)
	}
	return nil
}

func main() {
	err := boo()

	fmt.Println(err.Error()) // Which will return the string

}
