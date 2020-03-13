package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Array []interface{}
type Struct map[string]interface{}

var xmlSpecial = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
	'\'': "&apos;",
	'&':  "&amp;",
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if s, ok := xmlSpecial[c]; ok {
			b.WriteString(s)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

type valueNode struct {
	Type string `xml:"attr"`
	Body string `xml:"chardata"`
}

func next(p *xml.Decoder) (xml.Name, interface{}, error) {
	se, e := nextStart(p)
	if e != nil {
		return xml.Name{}, nil, e
	}

	var nv interface{}
	switch se.Name.Local {
	case "string":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		return xml.Name{}, s, nil
	case "boolean":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		s = strings.TrimSpace(s)
		var b bool
		switch s {
		case "true", "1":
			b = true
		case "false", "0":
			b = false
		default:
			e = errors.New("invalid boolean value")
		}
		return xml.Name{}, b, e
	case "int", "i1", "i2", "i4", "i8":
		var s string
		var i int
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		i, e = strconv.Atoi(strings.TrimSpace(s))
		return xml.Name{}, i, e
	case "double":
		var s string
		var f float64
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		f, e = strconv.ParseFloat(strings.TrimSpace(s), 64)
		return xml.Name{}, f, e
	case "dateTime.iso8601":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		t, e := time.Parse("20060102T15:04:05", s)
		if e != nil {
			t, e = time.Parse("2006-01-02T15:04:05-07:00", s)
			if e != nil {
				t, e = time.Parse("2006-01-02T15:04:05", s)
			}
		}
		return xml.Name{}, t, e
	case "base64":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		if b, e := base64.StdEncoding.DecodeString(s); e != nil {
			return xml.Name{}, nil, e
		} else {
			return xml.Name{}, b, nil
		}
	case "member":
		nextStart(p)
		return next(p)
	case "value":
		nextStart(p)
		return next(p)
	case "name":
		nextStart(p)
		return next(p)
	case "struct":
		st := Struct{}

		se, e = nextStart(p)
		for e == nil && se.Name.Local == "member" {
			// name
			se, e = nextStart(p)
			if se.Name.Local != "name" {
				return xml.Name{}, nil, errors.New("invalid response")
			}
			if e != nil {
				break
			}
			var name string
			if e = p.DecodeElement(&name, &se); e != nil {
				return xml.Name{}, nil, e
			}
			se, e = nextStart(p)
			if e != nil {
				break
			}

			// value
			_, value, e := next(p)
			if se.Name.Local != "value" {
				return xml.Name{}, nil, errors.New("invalid response")
			}
			if e != nil {
				break
			}
			st[name] = value

			se, e = nextStart(p)
			if e != nil {
				break
			}
		}
		return xml.Name{}, st, nil
	case "array":
		var ar Array
		nextStart(p) // data
		nextStart(p) // top of value
		for {
			_, value, e := next(p)
			if e != nil {
				break
			}
			ar = append(ar, value)

			if reflect.ValueOf(value).Kind() != reflect.Map {
				nextStart(p)
			}
		}
		return xml.Name{}, ar, nil
	case "nil":
		return xml.Name{}, nil, nil
	}

	if e = p.DecodeElement(nv, &se); e != nil {
		return xml.Name{}, nil, e
	}
	return se.Name, nv, e
}
func nextStart(p *xml.Decoder) (xml.StartElement, error) {
	for {
		t, e := p.Token()
		if e != nil {
			return xml.StartElement{}, e
		}
		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
	panic("unreachable")
}

func toXml(v interface{}, typ bool) (s string) {
	if v == nil {
		return "<nil/>"
	}
	r := reflect.ValueOf(v)
	t := r.Type()
	k := t.Kind()

	if b, ok := v.([]byte); ok {
		return "<base64>" + base64.StdEncoding.EncodeToString(b) + "</base64>"
	}

	switch k {
	case reflect.Invalid:
		panic("unsupported type")
	case reflect.Bool:
		return fmt.Sprintf("<boolean>%v</boolean>", v)
	case reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if typ {
			return fmt.Sprintf("<int>%v</int>", v)
		}
		return fmt.Sprintf("%v", v)
	case reflect.Uintptr:
		panic("unsupported type")
	case reflect.Float32, reflect.Float64:
		if typ {
			return fmt.Sprintf("<double>%v</double>", v)
		}
		return fmt.Sprintf("%v", v)
	case reflect.Complex64, reflect.Complex128:
		panic("unsupported type")
	case reflect.Array:
		s = "<array><data>"
		for n := 0; n < r.Len(); n++ {
			s += "<value>"
			s += toXml(r.Index(n).Interface(), typ)
			s += "</value>"
		}
		s += "</data></array>"
		return s
	case reflect.Chan:
		panic("unsupported type")
	case reflect.Func:
		panic("unsupported type")
	case reflect.Interface:
		return toXml(r.Elem(), typ)
	case reflect.Map:
		s = "<struct>"
		for _, key := range r.MapKeys() {
			s += "<member>"
			s += "<name>" + xmlEscape(key.Interface().(string)) + "</name>"
			s += "<value>" + toXml(r.MapIndex(key).Interface(), typ) + "</value>"
			s += "</member>"
		}
		s += "</struct>"
		return s
	case reflect.Ptr:
		panic("unsupported type")
	case reflect.Slice:
		s = "<array><data>"
		for n := 0; n < r.Len(); n++ {
			s += "<value>"
			s += toXml(r.Index(n).Interface(), typ)
			s += "</value>"
		}
		s += "</data></array>"
		return s
	case reflect.String:
		if typ {
			return fmt.Sprintf("<string>%v</string>", xmlEscape(v.(string)))
		}
		return xmlEscape(v.(string))
	case reflect.Struct:
		s = "<struct>"
		for n := 0; n < r.NumField(); n++ {
			s += "<member>"
			s += "<name>" + t.Field(n).Name + "</name>"
			s += "<value>" + toXml(r.FieldByIndex([]int{n}).Interface(), true) + "</value>"
			s += "</member>"
		}
		s += "</struct>"
		return s
	case reflect.UnsafePointer:
		return toXml(r.Elem(), typ)
	}
	return
}

// Client is client of XMLRPC
type Client struct {
	HttpClient *http.Client
	url        string
}

// NewClient create new Client
func NewClient(url string) *Client {
	return &Client{
		HttpClient: &http.Client{Transport: http.DefaultTransport, Timeout: 10 * time.Second},
		url:        url,
	}
}

func makeRequest(name string, args ...interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0"?><methodCall>`)
	buf.WriteString("<methodName>" + xmlEscape(name) + "</methodName>")
	buf.WriteString("<params>")
	for _, arg := range args {
		buf.WriteString("<param><value>")
		buf.WriteString(toXml(arg, true))
		buf.WriteString("</value></param>")
	}
	buf.WriteString("</params></methodCall>")
	return buf
}

func call(client *http.Client, url, name string, args ...interface{}) (v interface{}, e error) {
	r, e := client.Post(url, "text/xml", makeRequest(name, args...))
	if e != nil {
		return nil, e
	}

	// Since we do not always read the entire body, discard the rest, which
	// allows the http transport to reuse the connection.
	defer io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()

	if r.StatusCode/100 != 2 {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	p := xml.NewDecoder(r.Body)
	se, e := nextStart(p) // methodResponse
	if se.Name.Local != "methodResponse" {
		return nil, errors.New("invalid response: missing methodResponse")
	}
	se, e = nextStart(p) // params
	if se.Name.Local != "params" {
		return nil, errors.New("invalid response: missing params")
	}
	se, e = nextStart(p) // param
	if se.Name.Local != "param" {
		return nil, errors.New("invalid response: missing param")
	}
	se, e = nextStart(p) // value
	if se.Name.Local != "value" {
		return nil, errors.New("invalid response: missing value")
	}
	_, v, e = next(p)
	return v, e
}

// Call call remote procedures function name with args
func (c *Client) Call(name string, args ...interface{}) (v interface{}, e error) {
	return call(c.HttpClient, c.url, name, args...)
}

// Global httpClient allows us to pool/reuse connections and not wastefully
// re-create transports for each request.
var httpClient = &http.Client{Transport: http.DefaultTransport, Timeout: 10 * time.Second}

// Call call remote procedures function name with args
func Call(url, name string, args ...interface{}) (v interface{}, e error) {
	return call(httpClient, url, name, args...)
}
