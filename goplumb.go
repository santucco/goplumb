

/*2:*/


//line goplumb.w:52

// Copyright (c) 2013 Alexander Sychev. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * The name of author may not be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package goplumb provides interface to plumber - interprocess messaging from Plan 9.
package goplumb

import(


/*6:*/


//line goplumb.w:186

"os"



/*:6*/



/*9:*/


//line goplumb.w:200

"strings"



/*:9*/



/*13:*/


//line goplumb.w:251

"fmt"



/*:13*/



/*17:*/


//line goplumb.w:304

"errors"
"io"



/*:17*/



/*22:*/


//line goplumb.w:394

"bytes"
"strconv"



/*:22*/


//line goplumb.w:84

)

type(


/*4:*/


//line goplumb.w:163

//Message desribes a plumber message.
Message struct{
Src string
Dst string
Wdir string
Type string
Attr[]Attr
Data[]byte
}



/*:4*/



/*5:*/


//line goplumb.w:176

//Attr is a description of an attribute of a plumber message.
Attr struct{
Name string
Value string
}



/*:5*/



/*7:*/


//line goplumb.w:189

Plumb struct{
f*os.File
}



/*:7*/


//line goplumb.w:88

)

var(


/*8:*/


//line goplumb.w:194

//PlumberDir is a default mount point of plumber.
PlumberDir string= "/mnt/plumb/"



/*:8*/


//line goplumb.w:92

)



/*:2*/



/*10:*/


//line goplumb.w:204

//Open opens a specified port with a specified omode.
//If the port begin with a slash, it is taken as a literal file name,
//otherwise it is a file name in the plumber file system at PlumberDir.
func Open(port string,omode int)(*Plumb,error){
if!strings.HasPrefix(port,"/"){
if!strings.HasSuffix(PlumberDir,"/"){
PlumberDir+= "/"
}
port= PlumberDir+port
}
var p Plumb
var err error
if p.f,err= os.OpenFile(port,omode,0600);err!=nil{
return nil,err
}
return&p,nil
}



/*:10*/



/*12:*/


//line goplumb.w:238

//Send sends a message.
func(this*Plumb)Send(message*Message)error{
if this==nil||this.f==nil||message==nil{
return os.ErrInvalid
}
b:=Pack(message)
_,err:=this.f.Write(b)
return err
}



/*:12*/



/*14:*/


//line goplumb.w:255

//Pack packs a message to []byte.
func Pack(message*Message)[]byte{
s:=fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%d\n",
message.Src,message.Dst,
message.Wdir,message.Type,
PackAttr(message.Attr),
len(message.Data))
b:=append([]byte{},[]byte(s)...)
return append(b,message.Data...)
}



/*:14*/



/*15:*/


//line goplumb.w:269

//PackAttr packs attr to string. If an attribute value contains a white space,
//a quote or an equal sign the value will be quoted.
func PackAttr(attr[]Attr)string{
var s string
for i,v:=range attr{
if i!=0{
s+= " "
}
if strings.ContainsAny(v.Value," '=\t"){
s+= fmt.Sprintf("%s='%s'",v.Name,strings.Replace(v.Value,"'","''",-1))
}else{
s+= fmt.Sprintf("%s=%s",v.Name,v.Value)
}
}
return s
}



/*:15*/



/*16:*/


//line goplumb.w:288

//SendText sends a text-only message; it assumes Type is text and Attr is empty.
func(this*Plumb)SendText(src string,dst string,wdir string,data string)error{
m:=&Message{
Src:src,
Dst:dst,
Wdir:wdir,
Type:"text",
Data:[]byte(data)}
return this.Send(m)
}



/*:16*/



/*18:*/


//line goplumb.w:309

//Recv returns a received message or an error.
func(this*Plumb)Recv()(*Message,error){
if this==nil||this.f==nil{
return nil,os.ErrInvalid
}
b:=make([]byte,8192)
n,err:=this.f.Read(b)
if err!=nil{
return nil,err
}
m,r:=UnpackPartial(b[:n])
if m!=nil{
return m,nil
}
if r==0{
return nil,errors.New("buffer too small")
}
if r> len(b)-n{
b1:=make([]byte,r+n)
copy(b1,b)
b= b1
}else{
b= b[:n+r]
}
_,err= io.ReadFull(this.f,b[n:])
if err!=nil{
return nil,err
}
m,r= UnpackPartial(b)
if m!=nil{
return m,nil
}
return nil,errors.New("buffer too small")
}



/*:18*/



/*21:*/


//line goplumb.w:385

//Unpack return unpacked message.
func Unpack(b[]byte)*Message{
m,_:=UnpackPartial(b)
return m
}



/*:21*/



/*23:*/


//line goplumb.w:399

//UnpackPartial helps to unpack messages splited in peaces.
//The first call to UnpackPartial for a given message must be sufficient to unpack
//the header; subsequent calls permit unpacking messages with long data sections.
//For each call, b contans the complete message received so far.
//If the message is complete, a pointer to the resulting message m will be returned,
//and a number of remainings bytes r will be set to 0.
//Otherwise m will be nil and r will be set to the number of bytes
//remaining to be read for this message
//to be complete (recall that the byte count is in the header).
//Those bytes should be read by the caller, placed at location b[r:],
//and the message unpacked again.
//If an error is encountered, m will be nil and r will be zero.
func UnpackPartial(b[]byte)(m*Message,r int){
bb:=bytes.Split(b,[]byte{'\n'})
if len(bb)<6{
return nil,0
}
m= &Message{
Src:string(bb[0]),
Dst:string(bb[1]),
Wdir:string(bb[2]),
Type:string(bb[3]),
Attr:UnpackAttr(string(bb[4]))}
n,err:=strconv.Atoi(string(bb[5]))
if err!=nil{
return nil,0
}
i:=0
for j:=0;j<6;j++{
i+= len(bb[j])+1
}
if r= n-(len(b)-i);r!=0{
return nil,r
}
m.Data= make([]byte,n)
copy(m.Data,b[i:i+n])
return m,0
}



/*:23*/



/*25:*/


//line goplumb.w:479

//UnpackAttr unpack the attributes from s
func UnpackAttr(s string)[]Attr{
var attrs[]Attr
for i:=0;i<len(s);{
var a Attr
for;i<len(s)&&s[i]!='=';i++{
a.Name+= s[i:i+1]
}
i++
if i==len(s){
break
}
if s[i]=='\''{
i++
for;i<len(s);i++{
if s[i]=='\''{
if i+1==len(s){
break
}
if s[i+1]!='\''{
break
}
i++
}
a.Value+= s[i:i+1]
}
i++
}else{
for;i<len(s)&&s[i]!=' ';i++{
a.Value+= s[i:i+1]
}

}
i++
attrs= append(attrs,a)
}
return attrs
}



/*:25*/



/*26:*/


//line goplumb.w:520

//Close closes a plumbing connection.
func(this*Plumb)Close(){
if this.f!=nil{
this.f.Close()
this.f= nil
}
}



/*:26*/


