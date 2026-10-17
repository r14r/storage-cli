package storage

import (
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "io"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "sort"
 "strings"
 "time"
)
const Version = "0.1.0"
type File struct { Path string `json:"path"`; Size int64 `json:"size"` }
type Target struct { Name string; Path string; Deletable bool }
func targets() []Target {
 home,_:=os.UserHomeDir()
 result:=[]Target{{"app-cache",filepath.Join(home,"Library/Caches"),true},{"system-cache","/Library/Caches",false},{"system-temp","/private/var/tmp",false}}
 if runtime.GOOS!="darwin" { result=[]Target{{"app-cache",filepath.Join(home,".cache"),true},{"system-temp","/tmp",false}} }
 if runtime.GOOS=="darwin" {
  if out,err:=exec.Command("brew","--cache").Output();err==nil {
   p:=strings.TrimSpace(string(out))
   if p!="" { result=append(result,Target{"homebrew",p,true}) }
  }
 }
 return result
}
func format(n int64)string { x:=float64(n);for _,u:=range []string{"B","KiB","MiB","GiB","TiB"} {if x<1024{return fmt.Sprintf("%.1f %s",x,u)};x/=1024};return fmt.Sprintf("%.1f PiB",x)}
func walk(root string,visit func(File))error{
 return filepath.WalkDir(root,func(p string,d os.DirEntry,err error)error{
  if err!=nil{return nil}
  if d.IsDir()||d.Type()&os.ModeSymlink!=0{return nil}
  if !d.Type().IsRegular(){return nil}
  info,e:=d.Info();if e==nil {visit(File{p,info.Size()})}
  return nil
 })
}
func list(root string,limit int,min int64)([]File,error){
 if limit<1{return nil,errors.New("--limit must be positive")}
 st,e:=os.Stat(root);if e!=nil{return nil,e};if !st.IsDir(){return nil,errors.New("path is not a directory")}
 var a []File
 e=walk(root,func(f File){if f.Size>=min {a=append(a,f)}})
 sort.Slice(a,func(i,j int)bool{return a[i].Size>a[j].Size})
 if len(a)>limit {a=a[:limit]};return a,e
}
func usage(w io.Writer,path string)error{
 out,err:=exec.Command("df","-h",path).CombinedOutput();if err!=nil{return fmt.Errorf("df: %s: %w",strings.TrimSpace(string(out)),err)}
 fmt.Fprint(w,string(out))
 if runtime.GOOS=="darwin"&&path=="/" {
  data:="/System/Volumes/Data"
  if _,err:=os.Stat(data);err==nil {
   b,e:=exec.Command("df","-h",data).CombinedOutput();if e==nil {fmt.Fprintln(w,"Data volume:");fmt.Fprint(w,string(b))}
  }
 }
 return nil
}
func Run(args []string,w io.Writer)error{
 if len(args)==0 {args=[]string{"help"}}
 switch args[0] {
 case "help","-h","--help":
  return help(w, append([]string{}, args[1:]...))
 case "version","--version":fmt.Fprintln(w,Version);return nil
 case "disks":
  if len(args)>1 && (args[1]=="--help" || args[1]=="-h") {return help(w,[]string{"disks"})}
  fs:=flag.NewFlagSet("disks",flag.ContinueOnError);path:=fs.String("path","/","filesystem path");apfs:=fs.Bool("apfs",false,"APFS container report");if e:=fs.Parse(args[1:]);e!=nil{return e}
  if *apfs {if runtime.GOOS!="darwin"{return errors.New("APFS is macOS-only")};out,e:=exec.Command("diskutil","apfs","list").CombinedOutput();fmt.Fprint(w,string(out));return e}
  return usage(w,*path)
 case "largest","folders":
  if len(args)>1 && (args[1]=="--help" || args[1]=="-h") {return help(w,[]string{args[0]})}
  if args[0]=="folders" {return runFolders(args[1:],w)}
  fs:=flag.NewFlagSet(args[0],flag.ContinueOnError);path:=fs.String("path",".","scan root");limit:=fs.Int("limit",25,"max entries");min:=fs.Int64("min-mb",1,"minimum size MiB (largest)");jsonOut:=fs.Bool("json",false,"JSON output");if e:=fs.Parse(args[1:]);e!=nil{return e}
  if *min<0{return errors.New("--min-mb must be non-negative")}
  if args[0]=="largest" {
   a,e:=list(*path,*limit,*min*1024*1024);if e!=nil{return e}
   if *jsonOut{return json.NewEncoder(w).Encode(a)}
   for _,f:=range a {fmt.Fprintf(w,"%12s  %s\n",format(f.Size),f.Path)};return nil
  }
  ents,e:=os.ReadDir(*path);if e!=nil{return e};var a []File
  for _,d:=range ents {if !d.IsDir()||d.Type()&os.ModeSymlink!=0{continue};p:=filepath.Join(*path,d.Name());var n int64;_ =walk(p,func(f File){n+=f.Size});a=append(a,File{p,n})}
  sort.Slice(a,func(i,j int)bool{return a[i].Size>a[j].Size});if *limit<1{return errors.New("--limit must be positive")};if len(a)>*limit{a=a[:*limit]}
  if *jsonOut{return json.NewEncoder(w).Encode(a)}
  for _,f:=range a{fmt.Fprintf(w,"%12s  %s\n",format(f.Size),f.Path)};return nil
 case "targets":
  for _,t:=range targets(){status:="report-only";if t.Deletable{status="eligible for explicit cleanup"};fmt.Fprintf(w,"%-15s %-30s %s\n",t.Name,status,t.Path)};return nil
 case "cleanup":
  if len(args)>1 && (args[1]=="--help" || args[1]=="-h") {return help(w,[]string{"cleanup"})}
  fs:=flag.NewFlagSet("cleanup",flag.ContinueOnError);category:=fs.String("category","all","app-cache|homebrew|system-cache|system-temp|all");days:=fs.Int("older-than",7,"minimum file age in days");limit:=fs.Int("limit",40,"preview maximum");deleteFiles:=fs.Bool("delete",false,"delete eligible files");yes:=fs.Bool("yes",false,"explicit deletion confirmation");jsonOut:=fs.Bool("json",false,"JSON preview");if e:=fs.Parse(args[1:]);e!=nil{return e}
  if *days<0||*limit<1{return errors.New("invalid age or limit")}
  if *deleteFiles&&(!*yes||*category=="all"){return errors.New("deletion requires --yes and an explicit category")}
  if *deleteFiles&&*jsonOut{return errors.New("--json cannot be combined with deletion")}
  found:=false;var a []File
  for _,t:=range targets(){
   if *category!="all"&&*category!=t.Name{continue};found=true
   if *deleteFiles&&!t.Deletable{return errors.New("system cleanup is report-only")}
   if !t.Deletable&&*category=="all"{continue}
   root,e:=filepath.EvalSymlinks(t.Path);if e!=nil{continue}
   st,e:=os.Stat(root);if e!=nil||!st.IsDir(){continue}
   _=filepath.WalkDir(root,func(p string,d os.DirEntry,err error)error{
    if err!=nil||p==root{return nil}
    if d.Type()&os.ModeSymlink!=0{return nil}
    if d.IsDir(){return nil}
    if !d.Type().IsRegular(){return nil}
    info,e:=d.Info();if e!=nil||time.Since(info.ModTime())<time.Duration(*days)*24*time.Hour{return nil}
    a=append(a,File{p,info.Size()});return nil
   })
  }
  if !found{return errors.New("unknown cleanup category")}
  sort.Slice(a,func(i,j int)bool{return a[i].Size>a[j].Size})
  if *jsonOut{return json.NewEncoder(w).Encode(a)}
  var total int64;for _,f:=range a{total+=f.Size}
  fmt.Fprintf(w,"%d eligible files; potential space %s\n",len(a),format(total))
  n:=len(a);if n>*limit{n=*limit}
  for _,f:=range a[:n] {fmt.Fprintf(w,"%12s  %s\n",format(f.Size),f.Path)}
  if *deleteFiles {
   // Resolve the selected cache root again and prohibit directory/symlink deletion.
   for _,f:=range a{
    info,e:=os.Lstat(f.Path);if e!=nil||!info.Mode().IsRegular(){continue}
    if e:=os.Remove(f.Path);e!=nil{fmt.Fprintf(w,"SKIPPED %s: %v\n",f.Path,e)}
   }
  }else{fmt.Fprintln(w,"Preview only. For user cache cleanup: --category app-cache|homebrew --delete --yes")}
  return nil
 default:return fmt.Errorf("unknown command %q",args[0])
 }
}
