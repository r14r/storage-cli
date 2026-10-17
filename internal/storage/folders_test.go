package storage

import (
 "bytes"
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestFoldersDepthAndTotals(t *testing.T){
 root:=t.TempDir()
 child:=filepath.Join(root,"dir");nested:=filepath.Join(child,"nested")
 if err:=os.MkdirAll(nested,0700);err!=nil{t.Fatal(err)}
 for p,n:=range map[string]int{filepath.Join(root,"one"):5,filepath.Join(child,"two"):7,filepath.Join(nested,"three"):11}{
  if err:=os.WriteFile(p,make([]byte,n),0600);err!=nil{t.Fatal(err)}
 }
 for _,tc:=range []struct{depth string;expected int}{
  {"0",0},{"1",2},{"2",4},{"3",5},
 }{
  var b bytes.Buffer
  if err:=Run([]string{"folders",root,"--depth",tc.depth,"--json"},&b);err!=nil{t.Fatal(err)}
  var entries []folderEntry
  if err:=json.Unmarshal(b.Bytes(),&entries);err!=nil{t.Fatal(err)}
  if len(entries)!=tc.expected{t.Fatalf("depth %s got %d entries, want %d",tc.depth,len(entries),tc.expected)}
  for _,e:=range entries{if e.Path==child&&e.Size!=18{t.Fatalf("aggregate size got %d",e.Size)}}
 }
}
func TestFoldersFlagsAndHelp(t *testing.T){
 root:=t.TempDir()
 for _,args:=range [][]string{{"folders",root,"--depth","-1"},{"folders",root,"--depth","abc"},{"folders",root,root},{"folders",root,"--path",root}}{
  if err:=Run(args,&bytes.Buffer{});err==nil{t.Fatalf("expected error for %v",args)}
 }
 for _,args:=range [][]string{{"help"},{"help","folders"},{"folders","--help"}}{
  var b bytes.Buffer;if err:=Run(args,&b);err!=nil{t.Fatal(err)}
  if !strings.Contains(b.String(),"folders"){t.Fatalf("help lacks folders: %s",b.String())}
 }
}
