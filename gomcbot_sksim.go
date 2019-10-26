package main

import (
	"fmt"
	"log"
	"strings"
	"regexp"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	bot "github.com/Tnze/gomcbot"
	auth "github.com/Tnze/gomcbot/authenticate"
)

var dynmap *regexp.Regexp = regexp.MustCompile(`(?i)^(マップ|まっぷ|map|mappu|まp)$`)
var radio *regexp.Regexp = regexp.MustCompile(`(?i)^(ラジオ|らじお|radio|ラヂオ|らぢお|razio|rajio|れでぃお|レディオ|れいでぃお|レイディオ|redhio|reidhio)$`)
var wiki *regexp.Regexp = regexp.MustCompile(`(?i)^(うぃき|ウィキ|wiki)$`)
var sureddo *regexp.Regexp = regexp.MustCompile(`(?i)^(スレ|すれ|sure|スレッド|すれっど|sureddo|thread)$`)
var zishin *regexp.Regexp = regexp.MustCompile(`(?i)^((地震)|(自身)|(じ|ジ|ｼﾞ|ji|zi)(し|シ|ｼ|si|shi)(ん|ン|ﾝ|n|nn))$`)
var tokei *regexp.Regexp = regexp.MustCompile(`(?i)^(時計|とけい|トケイ|tokei|くろっく|クロック|clock|kurokku)$`)
var mahjan *regexp.Regexp = regexp.MustCompile(`(?i)^(麻雀|まーじゃん|マージャン|majan|ma-jan|じゃんまー)$`)
var duel *regexp.Regexp = regexp.MustCompile(`(?i)^(決闘|デュエル|でゅえる|duel|づえl)$`)
var pikuto *regexp.Regexp = regexp.MustCompile(`(?i)^(ぴくとせんす|ピクトセンス|pictsense|pikutosensu)$`)
var quiz *regexp.Regexp = regexp.MustCompile(`(?i)^(くいず|クイズ|kuizu|quiz|qmaclone|qma)$`)
var taihuu *regexp.Regexp = regexp.MustCompile(`(?i)^(台風|たいふう|タイフウ|taihu|taifu|taihuu|taifuu|typhoon|たいふーん|タイフーン|taihu-n|taifu-n|taihu-nn|taifu-nn)$`)
var yururowa *regexp.Regexp = regexp.MustCompile(`(?i)^(ゆ|ユ|ゅ|ュ|yu|湯)(る|ル|ru)(ろ|ロ|ro)(わ|ワ|ゎ|ヮ|wa)(！|!)?$`)
var kaigi *regexp.Regexp = regexp.MustCompile(`(?i)^((((け|ケ|ｹ|ke)(ん|ン|ﾝ|n|nn)|嫌)((も|モ|ﾓ|mo)(う|ウ|ｳ|u)*)|儲)*((か|カ|ｶ|ka)(い|イ|ｲ|i)|会)(議|ぎ|ギ|ｷﾞ|gi)|(議|ぎ|ギ|ｷﾞ|gi)((だ|ダ|ﾀﾞ|da)(い|イ|ｲ|i)|題))$`)
var wikievent *regexp.Regexp = regexp.MustCompile(`(?i)^([いイｲ][べベﾍﾞ][んンﾝ][とトﾄ]|(event)|(絵ヴぇんt)|(ibe(n|nn)to))(!|！)*$`)
var image *regexp.Regexp = regexp.MustCompile(`(?i)^(image|今げ)$`)
var image2 *regexp.Regexp = regexp.MustCompile(`(?i)^今げ`)

type Data struct {
    Id string `json:"id"`
    Password string `json:"password"`
    Server string `json:"server"`
    Port int `json:"port"`
}

func getWikiEvent() (ev string, err error){
	var URL string = "http://kenmomine.wiki.fc2.com/wiki/イベント企画"
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return "", fmt.Errorf("error")
	}
	table := doc.Find("table.table")
	t := strings.Split(table.Text(), "\n")
	t1 := strings.Split(t[2], "  ")[4]
	t2 := strings.Split(t[3], "  ")[4]
	ev = t1 + "   " + t2
	return ev, nil
}

func short_url(URL string) (res string, err error){
	URL = "https://is.gd/create.php?format=simple&url=" + URL
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return "", fmt.Errorf("error")
	}
	res = doc.Text()
	return res, nil
}

func google_search(msg []string) (URL string, err error){
	word := strings.Join(msg, " ")
	tmp := strings.Split(word, "(")
	if len(tmp) >= 2 {
		tmp = tmp[:len(tmp)-1] 
		word = strings.Join(tmp, "(")
	}
	URL = "https://www.google.co.jp/search?hl=ja&source=hp&q=" + url.QueryEscape(word)
	URL = url.QueryEscape(URL)
	URL, err = short_url(URL)
	if err != nil {
		return "", fmt.Errorf("error")
	}
	return URL, nil
}

func google_image(msg []string) (URL string, err error){
	word := strings.Join(msg, " ")
	tmp := strings.Split(word, "(")
	if len(tmp) >= 2 {
		tmp = tmp[:len(tmp)-1] 
		word = strings.Join(tmp, "(")
	}
	URL = "https://www.google.co.jp/search?hl=ja&source=hp&tbm=isch&q=" + url.QueryEscape(word)
	URL = url.QueryEscape(URL)
	URL, err = short_url(URL)
	if err != nil {
		return "", fmt.Errorf("error")
	}
	return URL, nil
}

func analyze_chat(game *bot.Game, txt string) (err error) {
	tmp1 := strings.Split(strings.Split(txt, ">")[0], "<")
	if len(tmp1) < 2 {
		return fmt.Errorf("txt is not message")
	} 
	name := tmp1[1]
	msg := strings.Split(txt, " ")
	if len(msg) < 2 {
		return fmt.Errorf("txt is not message")
	} 
	msg = msg[1:]
	if name == "sksim" {
		return fmt.Errorf("sksim is me")
	}
	if name == "Super_AI" {
		if len(msg) < 2 {
			return fmt.Errorf("txt is not message")
		} 
		name = msg[0]
		l := len(name)
		name = name[1:l-1]
		msg = msg[1:]
	}
	if len(msg) >= 2 && msg[0] == "sksim" {
		name_called(game, name, msg[1:])
	} else {
		chat_func(game, name, msg)
	}

	return nil
}

func chat_func(game *bot.Game, name string, msg []string) {
	if dynmap.MatchString(msg[0]) {
		game.Chat("http://kenmomine.club:8123/")
	} 
	if radio.MatchString(msg[0]) {
		game.Chat("https://cytube.xyz/r/kenmomine")
	} 
	if wiki.MatchString(msg[0]) {
		game.Chat("http://kenmomine.wiki.fc2.com/")
	} 
	if sureddo.MatchString(msg[0]) {
		game.Chat("https://ff2ch.syoboi.jp/?q=%E5%AB%8C%E5%84%B2minecraft%E9%83%A8")
	} 
	if zishin.MatchString(msg[0]) {
		game.Chat("http://www.kmoni.bosai.go.jp/new/")
	} 
	if tokei.MatchString(msg[0]) {
		game.Chat("https://www.nict.go.jp/JST/JST5.html")
	} 
	if mahjan.MatchString(msg[0]) {
		game.Chat("http://tenhou.net/make_lobby.html")
	} 
	if duel.MatchString(msg[0]) {
		game.Chat("https://godfield.net/")
	} 
	if pikuto.MatchString(msg[0]) {
		game.Chat("http://pictsense.com/")
	} 
	if quiz.MatchString(msg[0]) {
		game.Chat("http://kishibe.dyndns.tv/QMAClone/")
	} 
	if taihuu.MatchString(msg[0]) {
		game.Chat("https://tenki.jp/bousai/typhoon/")
	} 
	if yururowa.MatchString(msg[0]) {
		game.Chat("http://kenmomine.wiki.fc2.com/wiki/ゆるろわ")
	} 
	if kaigi.MatchString(msg[0]) {
		game.Chat("http://kenmomine.wiki.fc2.com/wiki/嫌儲会議")
	} 
	if (image.MatchString(msg[0]) || image2.MatchString(msg[0])) && len(msg) >= 2 {
		tmp := strings.Split(msg[0], "今げ")
		fmt.Println(msg)
		fmt.Println(tmp)
		fmt.Println(len(tmp))
		if len(tmp) >= 2 {
			msg[0] = strings.Join(tmp[1:], "今げ")
		} else {
			msg = msg[1:]
		}
		fmt.Println(msg)
		res, err := google_image(msg)
		if err == nil {
			game.Chat(res)
		} else {
			game.Chat("なんかエラーだって")
		}
	} 
}

func name_called(game *bot.Game, name string, msg []string) {
	if wikievent.MatchString(msg[0]) {
		ev, err := getWikiEvent()
		if err == nil {
			game.Chat(ev)
		} else {
			game.Chat("なんかエラーだって")
		}
		
	} else {
		res, err := google_search(msg)
		if err == nil {
			game.Chat(res)
		} else {
			game.Chat("なんかエラーだって")
		}
	} 
}

func main() {
	datafile := "data.json"
	f, err := ioutil.ReadFile(datafile)
    if err != nil {
		panic("error")
    }
	data := new(Data)
	err = json.Unmarshal(f, data)
	if err != nil {
        panic("error")
	}
	resp, err := auth.Authenticate(data.Id, data.Password)
	
	if err != nil {
		panic(err)
	}
	Auth := resp.ToAuth()

	//Join server
	game, err := Auth.JoinServer(data.Server, data.Port)
	if err != nil {
		panic(err)
	}

	//Handle game
	events := game.GetEvents()
	go game.HandleGame()
	
	if err != nil {
		log.Panic(err)
	}

	for e := range events {//Reciving events
		switch e.(type) {
		case bot.PlayerSpawnEvent:
			fmt.Println("ログイン成功")
		case bot.ChatMessageEvent: //chat message
			fmt.Println(e.(bot.ChatMessageEvent).Msg)
			var txt string
			for _, v := range e.(bot.ChatMessageEvent).Msg.Extra { txt += v.Text }
			go analyze_chat(game, txt)
		}

	}
}