package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"time"
	"strings"
)

var DB, dberr = sql.Open("sqlite3", "/home/aron/.Skype/noszalyaron4/main.db")

func init() {
	if dberr != nil {
		panic( dberr )
	}
}

type Message struct {
	Id int
	Author string
	Timestamp int64
	Body_xml []byte
}

var Counter map[string]int=make(map[string]int)
var links = 0
var LinkCounter map[string]int=make(map[string]int)
//var RuneMap = make(map[rune]int)

func main() {
	ret := make([]*Message, 0)
	
	rows, err := DB.Query("SELECT id, author,timestamp,body_xml FROM Messages WHERE chatname=?", "#noszalyaron4/$175379b26a0d7b45")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		s := Message{}
		if err := rows.Scan(&s.Id, &s.Author, &s.Timestamp, &s.Body_xml); err != nil {
			panic(err)
		}
		ret = append(ret, &s)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	
	
	maxlen := 0
	maxauthor := ""
	
	minlen := 1000000
	minauthor := ""

	t1 := int64(0)
	author1 := ""
	tmp1 := time.Time{}
	
	t2 := int64(10000000)
	author2 := ""
	tmp2 := time.Time{}
		
	longestauthor := ""	
	longest := 0
	
	currauthor := ""
	curr := 0

	var freq [24]int
	var freq_month [13]int
	
	ego := 0
	ilostthegame := 0
	enekora := 0
	
	bef := ""
	for _, s:=range ret {
		if bef == s.Author {
			curr++
			currauthor = s.Author
		}else {
			if longest <= curr {
				longest = curr
				longestauthor = currauthor
			}
			curr=1
			currauthor=s.Author
		}
		Counter[s.Author]++
		if len(s.Body_xml)>=maxlen {
			maxlen = len(s.Body_xml)
			maxauthor = s.Author 
		}
		if len(s.Body_xml)<=minlen {
			minlen = len(s.Body_xml)
			minauthor = s.Author 
		}
		t := time.Unix(s.Timestamp,0)
		freq[t.Hour()]++
		freq_month[t.Month()]++
		curr := int64(t.Hour()*3600+t.Minute()*60+t.Second())
		if curr >= t1 {
			t1 = curr
			author1 = s.Author
			tmp1 = t
		}
		if curr <= t2 {
			t2 = curr
			author2 = s.Author
			tmp2 = t
		}
		if strings.Contains(string(s.Body_xml),"http://") {
			links++
			LinkCounter[s.Author]++
		}
		if strings.Contains(strings.ToLower(string(s.Body_xml)),"matekmaffia") {
			ego++
		}
		if strings.Contains(strings.ToLower(string(s.Body_xml)),"vesztettem") {
			ilostthegame++
		}
		if strings.Contains(strings.ToLower(string(s.Body_xml))," ének") {
			enekora++
		}
		/*for _, r := range string(s.Body_xml) {
			RuneMap[r]++;
		}*/
		bef = s.Author
	}
	
	
	fmt.Println(time.Now().Format("2006 Jan 2 15:04:05"),"-i jelentése matekmaffia(r) chatszoba állapotáról\n")
	fmt.Println("Hozzászólások száma a matekmaffia(r) szobában:",len(ret))
	for i, s:=range Counter {
		fmt.Println(">",i, s, "darabot szólt hozzá")
	}
	fmt.Println("")
	fmt.Println("A leghosszabb hozzászólást:",maxauthor, "írta, ez",maxlen, "karakter hosszú volt! (Lehetséges hogy több darab is van, de csak az utolsót nézzük!)")
	fmt.Println("A legrövidebb hozzászólást:",minauthor, "írta, ez",minlen, "karakter hosszú volt! (Lehetséges hogy több darab is van, de csak az utolsót nézzük!)")
	fmt.Println("")
	fmt.Println("A legkorábbi hozzászólást:",author2, "írta, ez",tmp2.Format("2006 Jan 2 15:04:05"), "-kor történt! (Lehetséges hogy több darab is van, de csak az utolsót nézzük!)")
	fmt.Println("A legkésőbbi hozzászólást:",author1, "írta, ez",tmp1.Format("2006 Jan 2 15:04:05"), "-kor történt! (Lehetséges hogy több darab is van, de csak az utolsót nézzük!)")
	fmt.Println("")
	fmt.Println("A leghosszabb önmagával tartott diskurációt",longestauthor,"folytatta mintegy",longest, "hozzászóláson keresztül! (Lehetséges hogy több darab is van, de csak az utolsót nézzük!)")
	fmt.Println("")
	for i, l:= range freq {
		fmt.Println(">",i, ": 00-tól ",i,": 59-ig ",l,"hozzászólás van.")
		
	}
	fmt.Println("")
	for i, l:= range freq_month {
		fmt.Println(">",i, ". hónapban ",l,"hozzászólás van.")
	}
	fmt.Println("")
	fmt.Println("Öszessen",links,"db link lett a szobában posztolva!")
	fmt.Println("Felhasználók szerinti megoszlás:")
	for i, l:= range LinkCounter {
		fmt.Println(">",i,l)
	}
	fmt.Println("")
	/*for i, l:= range RuneMap {
		var s string = string(i)
		if i=='\n' {
			s="\\n"
		}
		fmt.Println(">",s, " karakterből a tagok ",l,"dbot írtak.")
	}*/
	fmt.Println("")
	
	fmt.Println("A MatekMaffia szó",ego,"-szer lett leírva!")
	fmt.Println("A MatekMaffia tagjai",ilostthegame,"-szer vesztették el A Játékot!")
	fmt.Println("A MatekMaffia tagjai",enekora,"-szer gondoltak az ének órára!")
	
}

