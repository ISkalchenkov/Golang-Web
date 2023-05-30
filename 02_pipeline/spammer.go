package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	wg := &sync.WaitGroup{}
	var in, out chan interface{}
	for _, c := range cmds {
		wg.Add(1)
		out = make(chan interface{})
		go func(in, out chan interface{}, command cmd) {
			defer wg.Done()
			command(in, out) // команда отработала - данные закончились,
			close(out)       // закрываем out, и следующая команда в пайплайне выйдет из цикла по каналу in
		}(in, out, c)
		in = out // вход следующей команды - выход предыдущей
	}
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	// 	in - string
	// 	out - User
	userAlreadyBeen := map[uint64]bool{}
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for email := range in {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			user := GetUser(email)
			mu.Lock()
			firstSelect := !userAlreadyBeen[user.ID]
			if firstSelect {
				userAlreadyBeen[user.ID] = true
			}
			mu.Unlock()
			if firstSelect {
				out <- user
			}
		}(email.(string))
	}
	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	// 	in - User
	// 	out - MsgID
	wg := &sync.WaitGroup{}
	batch := make([]User, 0, GetMessagesMaxUsersBatch)
	for user := range in {
		batch = append(batch, user.(User))
		if len(batch) != GetMessagesMaxUsersBatch {
			continue
		}
		wg.Add(1)
		go func(batch []User) {
			defer wg.Done()
			getAndSendMsgs(batch, out)
		}(batch)
		batch = make([]User, 0, GetMessagesMaxUsersBatch)
	}
	if len(batch) != 0 { // если input закрылся, а батч не заполнился до конца
		getAndSendMsgs(batch, out)
	}
	wg.Wait()
}

func getAndSendMsgs(batch []User, out chan interface{}) {
	msgIds, err := GetMessages(batch...)
	if err != nil {
		log.Println(err.Error())
	}
	for _, msgID := range msgIds {
		out <- msgID
	}
}

func CheckSpam(in, out chan interface{}) {
	// in - MsgID
	// out - MsgData
	wg := &sync.WaitGroup{}
	for i := 0; i < HasSpamMaxAsyncRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msgID := range in {
				hasSpam, err := HasSpam(msgID.(MsgID))
				if err != nil {
					log.Println(err.Error())
				}
				out <- MsgData{ID: msgID.(MsgID), HasSpam: hasSpam}
			}
		}()
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	// in - MsgData
	// out - string
	msgs := []MsgData{}
	for msgData := range in {
		msgs = append(msgs, msgData.(MsgData))
	}
	less := func(i, j int) bool {
		if msgs[i].HasSpam == msgs[j].HasSpam {
			return msgs[i].ID < msgs[j].ID
		}
		return msgs[i].HasSpam && !msgs[j].HasSpam
	}
	sort.Slice(msgs, less)
	for _, msgData := range msgs {
		out <- fmt.Sprintf("%t %d", msgData.HasSpam, msgData.ID)
	}
}
