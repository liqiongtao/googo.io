package goo_utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"
)

func TestStringMap_MarshalXML(t *testing.T) {
	data := StringMap{
		"user": "1",
	}

	var bf bytes.Buffer

	e := xml.NewEncoder(&bf)
	start := xml.StartElement{
		Name: xml.Name{Local: "xml"},
	}

	if err := data.MarshalXML(e, start); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bf.String())
}

func TestStringMap_UnMarshalXML(t *testing.T) {
	str := `
<xml>
  <ToUserName><![CDATA[toUser]]></ToUserName>
  <FromUserName><![CDATA[FromUser]]></FromUserName>
  <CreateTime>123456789</CreateTime>
  <MsgType><![CDATA[event]]></MsgType>
  <Event><![CDATA[subscribe]]></Event>
  <EventKey><![CDATA[qrscene_123123]]></EventKey>
  <Ticket><![CDATA[TICKET]]></Ticket>
</xml>
`

	var bf bytes.Buffer
	var data = &StringMap{}

	e := xml.NewDecoder(bytes.NewReader([]byte(str)))
	start := xml.StartElement{
		Name: xml.Name{Local: "xml"},
	}

	if err := data.UnmarshalXML(e, start); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bf.String())
	fmt.Println(data)
}
