package storage

import (
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "io"
 "os"
 "path/filepath"
 "sort"
)

type folderEntry struct {
 Path string `json:"path"`
 Size int64 `json:"size"`
 Kind string `json:"kind"`
 Depth int `json:"depth"`
}
func runFolders(args []string,w io.Writer) error {
 fs:=flag.NewFlagSet("folders",flag.ContinueOnError)
 fs.SetOutput(io.Discard)
 pathFlag:=fs.String("path","","directory to inspect")
 depth:=fs.Int("depth",1,"maximum nesting level, 1 = direct children")
 limit:=fs.Int("limit",0,"maximum displayed items; 0 = all")
 jsonOutput:=fs.Bool("json",false,"print JSON")
 // flag accepts flags before positional path; retain support for flags after it.
 var positional []string
 var flags []string
 for i:=0;i<len(args);i++{
  a:=args[i]
  if a=="--" {positional=append(positional,args[i+1:]...);break}
  if a=="--depth"||a=="--limit"||a=="--path"||a=="-depth"||a=="-limit"||a=="-path" {
   if i+1>=len(args){return fmt.Errorf("%s needs a value",a)}
   flags=append(flags,a,args[i+1]);i++;continue
  }
  if a=="--json"||a=="-json" {flags=append(flags,a);continue}
  if len(a)>0&&a[0]=='-' {flags=append(flags,a);continue}
  positional=append(positional,a)
 }
 if err:=fs.Parse(flags);err!=nil{return err}
 if len(positional)>1 {return errors.New("folders accepts one directory path")}
 if *pathFlag!=""&&len(positional)>0{return errors.New("use either positional FOLDER or --path")}
 root:=*pathFlag
 if len(positional)==1 {root=positional[0]}
 if root==""{root="."}
 if *depth<0||*limit<0{return errors.New("--depth and --limit must be non-negative")}
 info,err:=os.Stat(root);if err!=nil{return err};if !info.IsDir(){return errors.New("FOLDER must be a directory")}
 root,err=filepath.Abs(root);if err!=nil{return err}
 var result []folderEntry
 // Inspect each child; for directories compute recursive aggregate sizes once.
 var scan func(string,int)(int64,error)
 scan=func(dir string,level int)(int64,error){
  entries,err:=os.ReadDir(dir);if err!=nil{return 0,err}
  var sum int64
  for _,entry:=range entries {
   p:=filepath.Join(dir,entry.Name())
   if entry.Type()&os.ModeSymlink!=0 {continue}
   var n int64
   kind:="file"
   if entry.IsDir(){
    kind="folder"
    n,err=scan(p,level+1)
    if err!=nil{return 0,err}
   }else {
    if !entry.Type().IsRegular(){continue}
    stat,e:=entry.Info();if e!=nil{return 0,e};n=stat.Size()
   }
   sum+=n
   if level<=*depth {result=append(result,folderEntry{p,n,kind,level})}
  }
  return sum,nil
 }
 if _,err=scan(root,1);err!=nil{return err}
 sort.Slice(result,func(i,j int)bool{
  if result[i].Depth!=result[j].Depth{return result[i].Depth<result[j].Depth}
  if result[i].Size!=result[j].Size{return result[i].Size>result[j].Size}
  return result[i].Path<result[j].Path
 })
 if *limit>0&&len(result)>*limit{result=result[:*limit]}
 if *jsonOutput{return json.NewEncoder(w).Encode(result)}
 for _,v:=range result{
  fmt.Fprintf(w,"%12s  %s %s\n",format(v.Size),v.Kind, v.Path)
 }
 return nil
}
