package command

import (
	"fmt"
	"strings"
)

func Parse(input string) (Command, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return Command{}, fmt.Errorf("empty command")
	}

	parts := strings.Fields(input)

	return Command{
		Name: strings.ToUpper(parts[0]),
		Args: parts[1:],
	}, nil
}
