

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


//line goplumb.w:348

"errors"



/*:19*/



/*29:*/


//line goplumb.w:555

"time"


/*:29*/


//line goplumb.w:106

)

const rule= `type is text
src is Test
plumb to goplumb
`


func prepare(t*testing.T){
// checking for a running plumber instance
p,err:=os.Open(PlumberDir+"rules")
if err==nil{
t.Log("plumber started already")
p.Close()
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
defer f.Close()
_,err= f.Write([]byte(rule))
if err!=nil{
t.Fatal(err)
}
}

func compare(m1*Message,m2*Message)bool{
if m1.Src!=m2.Src||
m1.Dst!=m2.Dst||
m1.Wdir!=m2.Wdir||
m1.Type!=m2.Type||
len(m1.Attr)!=len(m2.Attr){
return false
}
for n,v:=range m1.Attr{
if m2.Attr[n]!=v{
return false
}
}
return bytes.Compare(m1.Data,m2.Data)==0
}



/*11:*/


//line goplumb.w:224

func TestOpen(t*testing.T){
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


//line goplumb.w:351

func TestSendRecv(t*testing.T){
rp,err:=Open("goplumb",os.O_RDONLY)
if err!=nil{
t.Fatal(err)
}
defer rp.Close()
sp,err:=Open("send",os.O_WRONLY)
if err!=nil{
t.Fatal(err)
}
defer sp.Close()
var m Message
m.Src= "Test"
m.Dst= "goplumb"
m.Wdir= "."
m.Type= "text"
m.Attr= make(Attrs)
m.Attr["attr1"]= "value1"
m.Attr["attr2"]= "value2"
m.Attr["attr3"]= "value = '3\t"
m.Data= []byte("1234567890")
if err:=sp.Send(&m);err!=nil{
t.Fatal(err)
}
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

func TestSendRecvBigMessage(t*testing.T){
rp,err:=Open("goplumb",os.O_RDONLY)
if err!=nil{
t.Fatal(err)
}
defer rp.Close()
sp,err:=Open("send",os.O_WRONLY)
if err!=nil{
t.Fatal(err)
}
defer sp.Close()
var m Message
m.Src= "Test"
m.Dst= "goplumb"
m.Wdir= "."
m.Type= "text"
m.Attr= make(Attrs)
m.Attr["attr1"]= "value1"
m.Attr["attr2"]= "value2"
m.Attr["attr3"]= "value = '3\t"
m.Data= make([]byte,0,9000)
for i:=0;i<900;i++{
m.Data= append(m.Data,[]byte("1234567890")...)
}
if err:=sp.Send(&m);err!=nil{
t.Fatal(err)
}
t.Logf("message %#v has been sent\n",m)
r,err:=rp.Recv()
t.Logf("message %#v has been received\n",*r)
if!compare(r,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:24*/



/*30:*/


//line goplumb.w:558

func TestMessageChannel(t*testing.T){
rp,err:=Open("goplumb",os.O_RDONLY)
if err!=nil{
t.Fatal(err)
}
defer rp.Close()
sp,err:=Open("send",os.O_WRONLY)
if err!=nil{
t.Fatal(err)
}
defer sp.Close()

var m Message
m.Src= "Test"
m.Dst= "goplumb"
m.Wdir= "."
m.Type= "text"
m.Attr= make(Attrs)
m.Attr["attr1"]= "value1"
m.Attr["attr2"]= "value2"
m.Attr["attr3"]= "value = '3\t"
m.Data= []byte("1234567890")
ch,err:=rp.MessageChannel()
if err!=nil{
t.Fatal(err)
}
if err:=sp.Send(&m);err!=nil{
t.Fatal(err)
}
t.Logf("message %#v has been sent\n",m)
<-time.NewTimer(time.Second).C
rm,ok:=<-ch
if!ok{
t.Fatal(errors.New("messages channel is closed"))
}
t.Logf("message %#v has been received\n",*rm)

if!compare(rm,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:30*/


//line goplumb.w:159





/*:3*/


