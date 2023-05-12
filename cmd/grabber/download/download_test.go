package download

// func Download_Test(t testing.T) {
// 	source := `d:\IN\IN_2022-05-11\tobot_s01_02_2010__hd.mp4`
// 	dest := `d:\IN\IN_2022-05-11\reports2.mp4`

// 	dj := StartNew(source, dest)
// 	//fmt.Println(err)
// 	// if err != nil {
// 	// 	fmt.Println(err.Error())
// 	// 	return
// 	// }

// 	draw_tick := time.NewTicker(3 * time.Second)
// 	done := false
// 	downloading := false

// 	i := 0
// 	for !done {
// 		i++
// 		select {
// 		case rs := <-dj.Listen():
// 			fmt.Printf("%v                       \n", rs.String())
// 			if rs.err != nil {
// 				fmt.Printf(rs.String())
// 				return
// 			}
// 			if rs.completed {
// 				return
// 			}
// 			if rs.terminated {
// 				return
// 			}
// 			if rs.progress > 1500000000 {
// 				dj.Kill()
// 			}
// 		case <-draw_tick.C:
// 			downloading = !downloading
// 			fmt.Println("")
// 			fmt.Println("downloading", downloading)
// 			if dj == nil {
// 				continue
// 			}
// 			switch downloading {
// 			case true:
// 				fmt.Println("GO CONTINUE")
// 				dj.Continue()
// 			case false:
// 				fmt.Println("GO PAUSE")
// 				dj.Pause()
// 			}

// 		}

// 	}
// }
