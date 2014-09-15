package main

import (
	"fmt"
	//"strings"
	//"os"
	//"io/ioutil"
	"runtime"
	"sync"
	nm "github.com/janinamass/needlemango"
	nms "github.com/janinamass/nmw_stdin"

)

func main(){
	var stdinContent string
	stdinContent = nms.Repeat()
	var allA, allB []nm.Sequence
	allA, allB = nms.SplitOnMarker(stdinContent,"")
	eblosum62 := nm.MakeSubstitutionMatrix("EBLOSUM62" )
	eblosum62.SetMap("EBLOSUM62")

	var maxCPU int = runtime.NumCPU()
	var wg sync.WaitGroup
	tasks :=make(chan nms.Task,maxCPU)
	var resultStr []string

	if (len(allA)*len(allB) < maxCPU){
		maxCPU = len(allA)*len(allB)
	}

	for i := 0; i < maxCPU; i++ {
        wg.Add(1)
        go func() {
            for t := range tasks {
				resultStr = append(resultStr,nm.PrettyPrint(nm.Nmw(t.A,t.B,t.Sm)))
            }
            wg.Done()
        }()
    }
	var seen map[string]bool
	seen = make(map[string]bool)
	var ok bool
	var tmpkey string
	for _,a := range allA{
		if a.GetSequence() ==""{continue}
		for _,b:=range allB{
			if b.GetSequence() ==""{continue}
			tmpkey = b.GetHeader()+a.GetHeader()
			ok = seen[tmpkey]
			if ok{continue}
			tmpkey = a.GetHeader()+b.GetHeader()
			ok = seen[tmpkey]
			if ok{continue}
			seen[a.GetHeader()+b.GetHeader()]=true
			seen[b.GetHeader()+a.GetHeader()]=true
			tasks <- nms.Task{A:a,B:b,Sm:eblosum62}

		}
	}
	close(tasks)
	wg.Wait()

	for _,v := range resultStr{
		fmt.Print(v)
	}
}
