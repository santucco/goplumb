

/*4:*/


//line goplumb.w:98

package goplumb

import(
"os/exec"
"testing"
"bytes"
"time"
"code.google.com/p/goplan9/plan9"
"code.google.com/p/goplan9/plan9/client"


/*23:*/


//line goplumb.w:381

"errors"



/*:23*/


//line goplumb.w:108

)

const rule= `type is text
src is Test
plumb to goplumb
`
var fs*client.Fsys



func prepare(t*testing.T){
// checking for a running plumber instance
var err error
fs,err= client.MountService("plumb")
if err==nil{
t.Log("plumber started already")
}else{
// start plumber
cmd:=exec.Command("plumber")
err= cmd.Run()
if err!=nil{
t.Fatal(err)
}
t.Log("plumber is starting, wait a second")
time.Sleep(time.Second)
}
fs,err= client.MountService("plumb")
if err!=nil{
t.Fatal(err)
}
// setting a rule for the test
f,err:=fs.Open("rules",plan9.OWRITE)
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



/*14:*/


//line goplumb.w:240

func TestOpen(t*testing.T){
prepare(t)
var err error
if sp,err= Open("send",plan9.OWRITE);err!=nil{
t.Fatal(err)
}
if rp,err= Open("goplumb",plan9.OREAD);err!=nil{
t.Fatal(err)
}
}



/*:14*/



/*24:*/


//line goplumb.w:384

func TestSendRecv(t*testing.T){
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
if err!=nil{
t.Fatal(err)
}
t.Logf("message %#v has been received\n",*r)
if!compare(r,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:24*/



/*28:*/


//line goplumb.w:467

func TestSendRecvBigMessage(t*testing.T){
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
if err!=nil{
t.Fatal(err)
}
t.Logf("message %#v has been received\n",*r)
if!compare(r,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:28*/



/*33:*/


//line goplumb.w:577

func TestMessageChannel(t*testing.T){
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
ch,err:=rp.MessageChannel(0)
if err!=nil{
t.Fatal(err)
}
if err:=sp.Send(&m);err!=nil{
t.Fatal(err)
}
t.Logf("message %#v has been sent\n",m)
time.Sleep(time.Second)
rm,ok:=<-ch
if!ok{
t.Fatal(errors.New("messages channel is closed"))
}
t.Logf("message %#v has been received\n",*rm)

if!compare(rm,&m){
t.Fatal(errors.New("messages is not matched"))
}
}



/*:33*/



/*34:*/


//line goplumb.w:610

func TestClose(t*testing.T){
rp.Close()
sp.Close()
}



/*:34*/


//line goplumb.w:167





/*:4*/


