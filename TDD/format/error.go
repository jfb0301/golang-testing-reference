package format 

import(
	"fmt"
	"errors"
)

func Error(expression, message string) string {
	return fmt.Errorf("[ERROR] Invalid expression: '%s' - '%s'", expression) 
}


