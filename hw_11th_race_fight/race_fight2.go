package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

var races = []string{"dog", "pig", "cat", "mouse", "horse", "bird", "fish", "bear", "rabbit", "tiger"}

func main() {
	// 做race表（種族：小孩數量）
	race_con := make(map[string]int)
	// 做context base
	ctx := context.Background()
	god_ctx, cancel := context.WithCancel(ctx)
	key := "race"
	// 10 races for 10 goroutine
	for _, race := range races {
		ancestor_ctx := context.WithValue(god_ctx, key, race)
		go ancestor_birth(ancestor_ctx)
		race_con[race] = 0
	}
	// 傳入造冊
	race_members <- race_con
	// 等待一秒停止生產
	time.Sleep(time.Second * 1)
	birth_end <- true
	// 等待一秒讓生育程式跑完
	time.Sleep(time.Second * 1)
	log.Println("戰鬥開始")
	go ancestor_fight()
	// go ancestor_change()
	go child_fight()
	// go child_change()

	time.Sleep(time.Second * 10)
	cancel()
}

func ancestor_fight() {
LOOP:
	for {
		time.Sleep(time.Millisecond * 500)
		select {
		// 有獲勝者時stop
		case <-stop_fight:
			break LOOP
		default:
			log.Println("祖先鬥")
			// 拿名冊
			all := <-race_members
			log.Println("祖先鬥拿冊")
			// 隨機戰死一名
			a_close <- true
			// 收到戰死者
			defeat := <-ancestor_defeat
			log.Println("收到戰死祖先身份")
			// 名冊剔除戰死祖先
			delete(all, defeat)
			// 蒐集剩餘種族
			remain_race := []string{}
			for k := range all {
				remain_race = append(remain_race, k)
			}
			log.Println("剩餘種族", remain_race)

			// 剩餘種族為1則宣布獲勝者
			if len(remain_race) == 1 {
				log.Println("獲勝者是", remain_race[0])
				stop_fight <- true
				race_members <- all
				break LOOP
			}

			// 還名冊
			race_members <- all
			log.Println("祖先鬥還冊")
		}
	}
}

func child_fight() {
LOOP:
	for {
		time.Sleep(time.Millisecond * 200)
		select {
		// 有獲勝者時stop
		case <-stop_fight:
			break LOOP
		default:
			log.Println("大亂鬥")
			// 拿名冊
			all := <-race_members
			log.Println("大亂鬥拿冊", all)
			for i := 1; i <= 150; i++ {
				select {
				// 隨機死亡150名小孩
				case c_close <- true:
				default:
					// 若block（代表無小孩可死）則做最後收屍
					for {
						// 最後一人獲勝
						if len(all) <= 1 {
							log.Println("獲勝者是", all)
							stop_fight <- true
							race_members <- all
							break LOOP
						}
						// 收屍
						defeat := <-child_defeat
						for k, v := range defeat {
							if _, ok := all[k]; ok {
								all[k] = all[k] - v
								if all[k] <= 0 {
									delete(all, k)
									ancestor_eliminate <- k
								}
							}
						}
					}
				}
			}
			// 暫存小孩屍體
			tmp := map[string]int{}
		LOOP2:
			for {
				select {
				// 收集小孩屍體
				case defeat := <-child_defeat:
					for k, v := range defeat {
						tmp[k] = tmp[k] + v
					}
					// 無其他屍體可收集
				case <-time.After(time.Millisecond * 1):
					// 暫存有屍體
					if len(tmp) != 0 {
						log.Println("開始收屍")
						log.Println("目前名冊", all)
						// 蒐集剩餘種族
						remain_race := []string{}
						for k, v := range tmp {
							// 確認此種族是否存在，避免時間差導致錯誤
							if _, ok := all[k]; ok {
								all[k] = all[k] - v
								if all[k] <= 0 {
									delete(all, k)
									// 告知祖先已無小孩
									ancestor_eliminate <- k
								}
							} else {
								// 剔除不存在種族
								delete(tmp, k)
							}
						}
						// 全死，要避免需逐個殺死
						if len(all) <= 0 {
							log.Println("全部死亡")
							stop_fight <- true
							break LOOP
						}

						// 陣亡小孩數
						log.Println("陣亡小孩", tmp)
						for k := range all {
							remain_race = append(remain_race, k)
						}
						log.Println("剩餘種族", remain_race)

						// 暫存清空
						tmp = map[string]int{}

						log.Println("收屍完冊更新", all)
						if len(all) == 1 {
							log.Println("獲勝者是", remain_race[0])
							stop_fight <- true
							break LOOP
						}

						select {
						// 若還有剩餘屍體（時間差導致沒蒐集到）
						case defeat := <-child_defeat:
							for k, v := range defeat {
								tmp[k] = tmp[k] + v
							}
							// 若無歸還名冊
						default:
							log.Println("大亂鬥還冊", all)
							race_members <- all
							break LOOP2
						}
						// 暫存無屍體則歸還名冊
					} else {
						race_members <- all
						break LOOP2
					}

				}

			}

		}

	}
}

func ancestor_birth(ctx context.Context) {
	key := "race"
	race_raw := ctx.Value(key)
	race := race_raw.(string)
	cancel_ctx, cancel := context.WithCancel(ctx)
	child_ctx := context.WithValue(cancel_ctx, key, race)
	//預計生幾個
	rand_num := rand.Intn(6)
	rand_num = rand_num + 5
	// 開生
	time.Sleep(time.Millisecond * 100)
	for i := 1; i <= rand_num; i++ {
		// 拿冊
		race_con := <-race_members
		go child_birth(child_ctx)
		// 造冊
		race_con[race] = race_con[race] + 1
		race_members <- race_con
	}
	log.Println("祖先生完", race)
	for {
		select {
		// 上帝賜死
		case <-ctx.Done():
			cancel()
			ancestor_defeat <- race
			log.Println("祖先:", race, "連帶死亡")
			return
			// 戰死
		case <-a_close:
			cancel()
			ancestor_defeat <- race
			log.Println("祖先:", race, "擊殺死亡")
			return
			// 無小孩悲憤死
		case which_race := <-ancestor_eliminate:
			if which_race == race {
				cancel()
				log.Println("祖先:", race, "自殺身亡")
				return
				// 非本族pass
			} else {
				ancestor_eliminate <- which_race
				time.Sleep(time.Millisecond * 1)
			}
		}
	}
}

func child_birth(ctx context.Context) {
	key := "race"
	race_raw := ctx.Value(key)
	race := race_raw.(string)
	child_ctx, cancel := context.WithCancel(ctx)
	//預計生幾個
	rand_num := rand.Intn(3)
	rand_num = rand_num + 1
	select {
	// 時間到停止生
	case <-birth_end:
		birth_end <- true
		select {
		case <-ctx.Done():
			cancel()
			defeat := map[string]int{race: 1}
			child_defeat <- defeat
			return
		case <-c_close:
			cancel()
			defeat := map[string]int{race: 1}
			child_defeat <- defeat
			return
		}
	default:
		time.Sleep(time.Millisecond * 100)
		// 開生
		for i := 1; i <= rand_num; i++ {
			// 拿冊
			race_con := <-race_members
			// 造冊
			go child_birth(child_ctx)
			race_con[race] = race_con[race] + 1
			race_members <- race_con
		}
		select {
		case <-ctx.Done():
			cancel()
			defeat := map[string]int{race: 1}
			child_defeat <- defeat
			return
		case <-c_close:
			cancel()
			defeat := map[string]int{race: 1}
			child_defeat <- defeat
			return
		}
	}
}
