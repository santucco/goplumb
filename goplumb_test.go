

/*3:*/


//line goplumb.w:97

package goplumb

import(
"os"
"os/exec"
"testing"
"bytes"
"syscall"


/*19:*/


//line goplumb.w:347

"errors"



/*:19*/


//line goplumb.w:106

)

const rule= `type is text
src is Test
plumb to goplumb
`


func prepare(t*testing.T){
// checking for a running plumber instance
_,err:=os.Open(PlumberDir+"rules")
if err==nil{
t.Log("plumber started already")
}else{
// start plumber
cmd:=exec.Command("plumber","-m",PlumberDir)
err= cmd.Run()
if err!=nil{
t.Fatal(err)
}
t.Log("plumber is starting, wait a second")
syscall.Nanosleep(&syscall.Timespec{Sec:1,},nil)
}
// setting a rule for the test
f,err:=os.OpenFile(PlumberDir+"rules",os.O_WRONLY,0600)
if err!=nil{
t.Fatal(err)
}
_,err= f.Write([]byte(rule))
if err!=nil{
t.Fatal(err)
}
f.Close()
}

func compare(m1*Message,m2*Message)bool{
if m1.Src!=m2.Src||
m1.Dst!=m2.Dst||
m1.Wdir!=m2.Wdir||
m1.Type!=m2.Type||
len(m1.Attr)!=len(m2.Attr){
return false
}
for i:=0;i<len(m1.Attr);i++{
if m1.Attr[i]!=m2.Attr[i]{
return false
}
}
return bytes.Compare(m1.Data,m2.Data)==0
}



/*11:*/


//line goplumb.w:225

func Test1(t*testing.T){
prepare(t)
if _,err:=Open("send",os.O_WRONLY);err!=nil{
t.Fatal(err)
}

if _,err:=Open("goplumb",os.O_RDONLY);err!=nil{
t.Fatal(err)
}
}



/*:11*/



/*20:*/


//line goplumb.w:350

func Test2(t*testing.T){
rp,err:=Open("goplumb",os.O_RDONLY)
if err!=nil{
t.Fatal(err)
}

sp,err:=Open("send",os.O_WRONLY)
if err!=nil{
t.Fatal(err)
}

var m Message
m.Src= "Test"
m.Dst= "goplumb"
m.Wdir= "."
m.Type= "text"
m.Attr= append([]Attr{},
Attr{Name:"attr1",Value:"value1"},
Attr{Name:"attr2",Value:"value2"},
Attr{Name:"attr3",Value:"value = '3\t"})
m.Data= []byte("1234567890")
if err:=sp.Send(&m);err!=nil{
t.Fatal(err)
}
sp.Close()
t.Logf("message %#v has been sent\n",m)
r,err:=rp.Recv()
t.Logf("message %#v has been received\n",*r)
if!compare(r,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:20*/



/*24:*/


//line goplumb.w:441

func Test3(t*testing.T){
rp,err:=Open("goplumb",os.O_RDONLY)
if err!=nil{
t.Fatal(err)
}

sp,err:=Open("send",os.O_WRONLY)
if err!=nil{
t.Fatal(err)
}

var m Message
m.Src= "Test"
m.Dst= "goplumb"
m.Wdir= "."
m.Type= "text"
m.Attr= append([]Attr{},
Attr{Name:"attr1",Value:"value1"},
Attr{Name:"attr2",Value:"value2"},
Attr{Name:"attr3",Value:"value = '3\t"})
m.Data= make([]byte,0,9000)
for i:=0;i<900;i++{
m.Data= append(m.Data,[]byte("1234567890")...)
}
if err:=sp.Send(&m);err!=nil{
t.Fatal(err)
}
sp.Close()
t.Logf("message %#v has been sent\n",m)
r,err:=rp.Recv()
t.Logf("message %#v has been received\n",*r)
if!compare(r,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:24*/


//line goplumb.w:158





/*:3*/


