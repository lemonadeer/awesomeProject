package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)


func testReflect(inter interface{})  {
	getType := reflect.TypeOf(inter)
	fmt.Println("getType:", getType)
	getValue := reflect.ValueOf(inter)
	fmt.Println("getValue:", getValue)

	value := getValue.Field(0)
	fmt.Println(value)
	m := getType.Method(0)

	fmt.Println(m.Name)
	//field := getType.Field(0)
	//fmt.Println(field.Name)
	//fmt.Println(getType.NumField())
}

type a struct {
	str string
}

type pbf struct {
	n []*a
}

func countLines(i io.Reader, countMap map[string]int)  {
	input := bufio.NewScanner(i)
	for input.Scan() {
		text := input.Text()
		if text == "q" {
			return
		}
		countMap[text] += 1
	}
}

func any(i interface{}) string {
	return formatAtom(reflect.ValueOf(i))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "Invalid"
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

type myStruct struct {
	Name []string
}

type Movie struct {
	Title string
	Year int `json:"released"`
	Color bool `json:"color,omitempty"`
	Actors []string
}

type Person struct {
	Name string	`json:"name"`
	Work string	`json:"work"`
	Age int	`json:"age"`
}

type Persons struct {
	Persons []Person 	`json:"persons"`
	City string			`json:"city"`
}

type PortRange struct {
	From uint32
	To   uint32
}

type InboundDetourConfig struct {
	Protocol  string     `json:"protocol"`
	PortRange *PortRange `json:"port"`
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

type MyMux struct {

}

func sayHello(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "hello world")
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		sayHello(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func login(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t,_ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, token))
	} else {
		if err := r.ParseForm(); err!=nil {
			fmt.Println(err.Error())
		}
		if len(r.Form.Get("username")) == 0 {
			fmt.Fprint(w, "username is empty")
		}
		if _, err := strconv.Atoi(r.Form.Get("age")); err != nil {
			fmt.Fprint(w, "Input invalid age")
		}

		token := r.Form.Get("token")
		if token != "" {

		} else {

		}

		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		template.HTMLEscape(w, []byte(r.Form.Get("username")))
	}
}

func upload(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t,_ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		if err := r.ParseMultipartForm(32 << 20); err!=nil {
			fmt.Println(err.Error())
		}
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		fmt.Fprint(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	/*
	fmt.Println(errors.New("eof") == errors.New("eof"))

	f,err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(f))

	in := InboundDetourConfig{}
	err = json.Unmarshal(f, &in)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(*in.PortRange)

	/*
	movies := []Movie{
		{
			Title:  "he",
			Year:   1999,
			Color:  false,
			Actors: []string{"r", "jasas"},
		},
		{
			Title:  "fkeja",
			Year:   1979,
			Color:  false,
			Actors: []string{"rzx", "jke"},
		},
	}
	
	bytes,_ := json.MarshalIndent(movies, "", "    ")
	fmt.Println(string(bytes))
	moviesBack := new([]Movie)
	json.Unmarshal(bytes, moviesBack)
	fmt.Println(moviesBack)

	/*
	counts := make(map[string]int)
	argsNum := os.Args[1: ]
	if len(argsNum) == 0 {
		countLines(os.Stdin, counts)
	}
	fmt.Println(counts)

*/


	/*
	f := pbf{n:[]*a{
		&a{
			str:"rzx",
		},
		&a{
			str:"ehkdf",
		},
	}}
	for _,v := range f.n {
		fmt.Println(*v)
	}

	p := pb.Person{
		Name:                 "rzx",
		Id:                   210,
		Email:                "1399072433@qq.com",
		Phones:               []*pb.Person_PhoneNumber{
									&pb.Person_PhoneNumber{
										Number:               "15541924241",
										Type:                 pb.Person_MOBILE,
									},
									&pb.Person_PhoneNumber{
										Number:               "1323523432",
										Type:                 pb.Person_HOME,
									},
								},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	adBook := &pb.AddressBook{
		People:               []*pb.Person{
			&p,
		},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	fmt.Println(adBook)
	out, err := proto.Marshal(adBook)
	if err != nil {
		fmt.Println("marshal fialed")
	}
	if err := ioutil.WriteFile("pb", out, 0644); err != nil {
		fmt.Println("error")
	}

	in, err := ioutil.ReadFile("pb")
	if err != nil {
		fmt.Println(err.Error())
	}

	newAddrbook := &pb.AddressBook{}
	if err := proto.Unmarshal(in, newAddrbook); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(newAddrbook)


*/






	/*
	fmt.Println("Starting server")

	listener,err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Error listener", err.Error())
	}

	for {
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println("Error accept", err.Error())
			return
		}
		go serverHandler(conn)

	}
*/
}

func serverHandler(conn net.Conn)  {

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error Reading", err.Error())
	}
	fmt.Println("Receive num:", n, string(buf[:len(buf)]))

}













