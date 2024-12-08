package assistant

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const instructions = `You are an assistant equipped with a powerful calculator tool. You can solve mathematical expressions provided by the user. When the user asks a math question, use the calculator tool to evaluate the expression and return the result.

Examples:
- User: What is 5 + 7?
- Assistant: 5 + 7 = 12

- User: Can you calculate (3 * (2 + 4)) / 3?
- Assistant: (3 * (2 + 4)) / 3 = 6`

func SetUp() {

}

func TestMain(m *testing.M) {
	SetUp()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestListDelers(t *testing.T) {

	assistant, err := NewAssistant(WithAssistantID(""))
	list, err := assistant.ListAssistants()

	for _, a := range list {
		fmt.Println(a.ID)
		_, err = a.DeleteAssistant()
		assert.NoError(t, err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, assistant)

}
