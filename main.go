package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type queue[T any] struct{
	Urls []T
}

func (q *queue[T]) enqueue(toQueue T) {
	q.Urls = append(q.Urls, toQueue)
}

func (q *queue[T]) dequeue() (T,error){
	var zero T
	if len(q.Urls)==0{
		return zero,fmt.Errorf("queue is empty")
	}
	fmt.Println(len(q.Urls))
	res:=q.Urls[0]
	q.Urls[0]=zero
	q.Urls=q.Urls[1:]
	return res,nil
}


func main() {
	urls:=[]string{"https://web-scraping.dev/"}
	htmlUrls:=map[string][]string{}

	for _,url:=range urls{
		link,err:=downloadLink(url)
		if err!=nil{
			fmt.Println(err)
			return
		}
		htmlUrls[url]=link
	}

	for key,value:=range htmlUrls{
		urlNormalization(key,value)
		fmt.Println(key,value)
	}
	
}


func downloadLink(url string)([]string,error){
	resp,err:=http.Get(url)
	if err!=nil{
		return nil,err
	}
	token:=html.NewTokenizer(resp.Body) 
	result:=[]string{}
	outer:
	for{
		tokenType:=token.Next()
		switch tokenType{
		case html.ErrorToken:
			return nil,fmt.Errorf("error: %v",token.Token().Data)
		case html.StartTagToken:
			t:=token.Token()
			if t.Data=="a"{

				for _,attribute:=range t.Attr{
					if attribute.Key=="href"{
						result=append(result,attribute.Val)
					}
				}
			}
		case html.EndTagToken:
			if token.Token().Data=="html"{
				break outer
			}
		}
	}
	return result,nil
}

func urlNormalization(baseUrl string,urls []string){
	for id,url :=range urls{
		if strings.Contains(url,"https"){
			continue
		}
		urls[id]=baseUrl+url
	}
}