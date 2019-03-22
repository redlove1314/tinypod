package lib

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
)

func Test1(t *testing.T) {
	f, _ := GetFile("C:\\Users\\hehety\\Downloads\\folder_document_233px_1185632_easyicon.net.png")
	bs, _ := ioutil.ReadAll(f)
	fmt.Print("[]byte{")
	for i, v := range bs {
		if i == len(bs)-1 {
			fmt.Print(v)
		} else {
			fmt.Print(v, ",")
		}
	}
	fmt.Println("}")
}
func Test2(t *testing.T) {
	a := []string{"1.jpg", "11.jpg", "2.jpg"}
	fmt.Println(a)
	fmt.Println(regexp.Match("^/[^:]+:.+$", []byte("///:app/:http://12312323.com/::asdiqweqd")))
}
