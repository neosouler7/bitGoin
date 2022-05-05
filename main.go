package main

import (
	"github.com/neosouler7/bitGoin/cli"
	"github.com/neosouler7/bitGoin/db"
)

// TODO.
// 0. sender, receiver 함수 분리로 시작
// 1. channel 생성/초기화 시점 : close(c)
// 2. ok를 통한 close 여부 확인 : a, ok := <- c
// 3. send only, receive only channel 명시화
// 4. buffered channels
func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
