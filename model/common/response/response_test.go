package response

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestAppendError(t *testing.T) {
	var strErr []string
	err2 := errors.New("err2")
	strErr = append(strErr, err2.Error())
	join := strings.Join(strErr, ",")
	fmt.Println(join)
}
