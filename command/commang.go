package command

import (
	"fmt"
	"io"
	"os"

	"github.com/kynefuk/goLGTM/imgsource"
)

// Command represents Command object.
type Command struct {
	OutStream, ErrStream io.Writer
}

// Run is a main logic of command
func (c *Command) Run(source, message string) error {
	// 画像ソースごとに構造体を定義する
	// 入力された画像ソースを元に画像ソースの構造体を返す関数
	// メッセージをのせて画像を生成するクラス

	// sourceの値をチェック&構造体を返す構造体

	factory := imgsource.NewImgSrcFactory(source)
	imgSrc := factory.GetImgSrc()
	err := imgSrc.AddMessage(message)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "error happened. error: %s\n", err)
		return err
	}
	return nil
}

// NewCommand is a constructor of Command
func NewCommand() *Command {
	return &Command{OutStream: os.Stdout, ErrStream: os.Stderr}
}
