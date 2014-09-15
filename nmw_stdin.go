package nmw_stdin

import (
	"fmt"
//"strconv"
	"strings"
	"os"
	"io/ioutil"
	//"runtime"
	//"sync"
	nm "github.com/janinamass/needlemango"
)

func Repeat() string {
	bytes, err := ioutil.ReadAll(os.Stdin)
	check(err)
	return(string(bytes))
}
func ReadFromStdin() []byte{
	bytes, err := ioutil.ReadAll(os.Stdin)
	check(err)
	return(bytes)
}


/*
MAIN
*/

/*func main(){
	var stdinContent string
	stdinContent = repeat()
	var allA, allB []nm.Sequence
	allA, allB = splitOnMarker(stdinContent,"")
	eblosum62 := nm.MakeSubstitutionMatrix("EBLOSUM62" )
	eblosum62.SetMap("EBLOSUM62")

	var maxCPU int = runtime.NumCPU()
	var wg sync.WaitGroup
	tasks :=make(chan Task,10000)
	var resultStr []string

	if (len(allA)*len(allB) < maxCPU){
		maxCPU = len(allA)*len(allB)
	}

	for i := 0; i < maxCPU; i++ {
        wg.Add(1)
        go func() {
            for t := range tasks {
				resultStr = append(resultStr,nm.PrettyPrint(nm.Nmw(t.a,t.b,t.sm)))
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
			tasks <- Task{a:a,b:b,sm:eblosum62}

		}
	}
	close(tasks)
	wg.Wait()

	for _,v := range resultStr{
		fmt.Print(v)
	}
}
*/
func Consumer(limit int, inChan <-chan Task){
	for i := 0; i< limit;i++{
		for s := range inChan {
			fmt.Println(nm.PrettyPrint(nm.Nmw(s.A, s.B, s.Sm)))
    	}
	}
}

func Producer(allA []nm.Sequence, allB []nm.Sequence, sm nm.SubstitutionMatrix) <-chan Task {
    ch := make(chan Task)
    go func() {
		for _,a :=range allA{
			for _,b:= range allB{
				ch <- Task{a,b,sm}
			}
        }
        close(ch)
    }()
    return ch
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func SplitOnMarker(text string, marker string)([]nm.Sequence, []nm.Sequence){
	var resA []nm.Sequence
	var resB []nm.Sequence
	var pieces []string

	if marker == "" {
		marker = "終"
	}
	pieces = strings.Split(text, marker) //split in two pieces
	tmpseq := strings.Split(pieces[0], ">")
	for _,w := range(tmpseq){
		var tmp []string
		tmp = strings.Split(w,"\n")
		shortHeader := strings.Split(tmp[0], " ")[0]
		resA = append(resA, nm.MakeSequence(shortHeader, strings.Join(tmp[1:], "")))
	}

	tmpseq = strings.Split(pieces[1], ">")
	for _,w := range(tmpseq){
		var tmp []string
		tmp = strings.Split(w,"\n")
		shortHeader := strings.Split(tmp[0], " ")[0]
		resB = append(resB, nm.MakeSequence(shortHeader, strings.Join(tmp[1:], "")))
	}

	return resA, resB
}

type Task struct{
	A nm.Sequence
	B nm.Sequence
	Sm nm.SubstitutionMatrix
}

