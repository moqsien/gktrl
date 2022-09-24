package goktrl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/abiosoft/ishell/v2"
	"github.com/gin-gonic/gin"
)

type ContextType int

const (
	ContextClient ContextType = 1 // 客户端
	ContextServer ContextType = 2 // 服务端
	ArgsFormatStr string      = "args%sargs"
)

type Context struct {
	*gin.Context
	ShellContext  *ishell.Context
	Type          ContextType
	Options       KtrlOpt
	Args          []string
	Parser        *ParserPlus
	Table         *KtrlTable
	KtrlPath      string
	Client        *KtrlClient
	DefaultSocket string
	ShellCmdName  string
	Result        []byte
}

func (that *Context) GetResult(sockName ...string) ([]byte, error) {
	sName := that.DefaultSocket
	if len(sockName) > 0 && len(sockName[0]) > 0 {
		sName = sockName[0]
	}
	params := that.Parser.Params
	params[fmt.Sprintf(ArgsFormatStr, that.ShellCmdName)] = strings.Join(that.Args, ",")
	return that.Client.GetResult(that.KtrlPath, params, sName)
}

func (that *Context) Send(content interface{}, code ...int) {
	statusCode := http.StatusOK
	if len(code) > 0 {
		statusCode = code[0]
	}
	switch content.(type) {
	case string:
		r, _ := content.(string)
		that.Context.String(statusCode, r)
	case []byte:
		r, _ := content.([]byte)
		that.Context.String(statusCode, string(r))
	default:
		r, err := json.Marshal(content)
		if err != nil {
			fmt.Println(err)
			that.Context.String(http.StatusInternalServerError, err.Error())
			return
		}
		that.Context.String(statusCode, string(r))
	}
}
