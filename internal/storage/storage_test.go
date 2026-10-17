package storage

import (
 "bytes"
 "os"
 "path/filepath"
 "strings"
 "testing"
)
func TestVersion(t *testing.T){var b bytes.Buffer;if err:=Run([]string{"version"},&b);err!=nil||strings.TrimSpace(b.String())!=Version{t.Fatalf("%q %v",b.String(),err)}}
func TestLargest(t *testing.T){d:=t.TempDir();if e:=os.WriteFile(filepath.Join(d,"big"),make([]byte,1024),0600);e!=nil{t.Fatal(e)};a,e:=list(d,5,0);if e!=nil||len(a)!=1||a[0].Size!=1024{t.Fatalf("%v %v",a,e)}}
func TestDeletionGuards(t *testing.T){for _,args:=range [][]string{{"cleanup","--delete","--yes"},{"cleanup","--category","app-cache","--delete"},{"cleanup","--category","system-temp","--delete","--yes"}}{if err:=Run(args,&bytes.Buffer{});err==nil{t.Fatalf("accepted %v",args)}}}
