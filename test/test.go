package main

//func takeKey(data ...any) string {
//	builder := strings.Builder{}
//	for i, v := range data {
//		switch v.(type) {
//		case string:
//			builder.WriteString(v.(string))
//		case int32:
//			builder.WriteString(strconv.FormatInt(int64(v.(int32)), 10))
//		case uint32:
//			builder.WriteString(strconv.FormatInt(int64(v.(uint32)), 10))
//		case int64:
//			builder.WriteString(strconv.FormatInt(v.(int64), 10))
//		case uint64:
//			builder.WriteString(strconv.FormatInt(int64(v.(uint64)), 10))
//		case int:
//			builder.WriteString(strconv.Itoa(v.(int)))
//		default:
//			builder.WriteString("???")
//		}
//		if i != len(data)-1 {
//			builder.WriteString(":")
//		}
//	}
//
//	return builder.String()
//}

//func init() {
//	log.SetFormatter(&log.JSONFormatter{
//		TimestampFormat: "2006-01-02",
//		FieldMap: log.FieldMap{
//			log.FieldKeyTime: "@timestamp",
//		},
//	})
//
//	log.SetFormatter(&log.TextFormatter{
//		FullTimestamp: true,
//		FieldMap: log.FieldMap{
//			log.FieldKeyTime:  "时间",
//			log.FieldKeyLevel: "日志类型",
//			log.FieldKeyMsg:   "日志内容",
//		},
//	})
//}

var publicKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8cRR/YqjlW4xZaO8Zwev\nvD34pk95V978PxVzQ9zsEjbXPTmVtoGP0Iat7Fy0MpybrKGmeLiuy7LeVr7KNyGz\ndYVsQSm7VAsuBAzHXkjKoWagYWWP6TIgv2Y3uBW3tU1xdhjC4wh+Qv6RajWL/Wj1\nqRfSjtyXgD+09cwKccyl8Oc8jNGKZzX2+QSK7LbPPW3hAPwypaTX4KOpvIzzU7t1\nONs1aCY8c3tpNGi3IP/fqIYagSYJ55POWBFeSm7f/CIDaRJ7vziCZpAzOTf7WTrg\nNV6gHJEwAtT+kB+k/L0EyjwPsTXLW47rdhTdqzpEVdY9Jr8/pM7sR1S2bjSawBwP\nsQIDAQAB\n-----END PUBLIC KEY-----"

var privateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA8cRR/YqjlW4xZaO8ZwevvD34pk95V978PxVzQ9zsEjbXPTmV\ntoGP0Iat7Fy0MpybrKGmeLiuy7LeVr7KNyGzdYVsQSm7VAsuBAzHXkjKoWagYWWP\n6TIgv2Y3uBW3tU1xdhjC4wh+Qv6RajWL/Wj1qRfSjtyXgD+09cwKccyl8Oc8jNGK\nZzX2+QSK7LbPPW3hAPwypaTX4KOpvIzzU7t1ONs1aCY8c3tpNGi3IP/fqIYagSYJ\n55POWBFeSm7f/CIDaRJ7vziCZpAzOTf7WTrgNV6gHJEwAtT+kB+k/L0EyjwPsTXL\nW47rdhTdqzpEVdY9Jr8/pM7sR1S2bjSawBwPsQIDAQABAoIBAB0UmzQfF/wibAio\nwEG4V/gRkDYY+ySJqte/scSo7zBlrlAr/Ake3nibqpHyuK4ZzlPegdKljEjuM/ZF\nLreg8yAgs1vHNEQwsBFGpDiAEveFC6eLetr27592IR+gZR+GuC4XXmHGpMFUM5ON\n60/I7zuupOIQQJzpjM/AAkWb8x3dxr5D+Y+5mrpdszRioowoJNrexeEipg/nFCjS\nFEx/xN52S+LpV6xIbnZ3r5lTMdK7V26E7UUpkOA3prbswDpMiSF21vcjM317QpJd\n+/Mmfk0wKmeaSbFd/vqz5oQ7ou/vkIjJHGYlV0vbc6MTpWQr1lDzK3ihk3AM4dT9\nmar13k0CgYEA+TRsj7wAhnWkk6Y9Arr7iqfUo3BBFf7ZffVOK8hbSgFZ0ffQrIkH\nq2fWk/sxN9SbtF6iDqBpY+OTtkH7dUSp2Dh1p8zGOx8wTXqmQyjyAi0uKWKC2ilJ\nKzWDned63VpZMDMcVJuNNfIzIJh3ojQTl/G6RDhWHGgV2RlCcavyre8CgYEA+Fv5\n0UXhrAIcVaUMS6feJ1tHCSerCI/iIrBHFtRYO6RCNp3bTrfinn8c/31mHSmDpJaz\nAxn/lIRi7oR1IhGQBhp/hYVTPpX42kr5C2KS7ZC7NcB++zkx7HeEHnXubZjJ0c8r\nH85myghrlCM/lOumT7UeCgLtQ9HbAghbIW3OvF8CgYEAucZzB9PHMHWS8t8CrH5n\n9r2WryCH5LXPvS6Zz9nU9B59ryFm1rhwlz8Zn8eqsUw1pwjFFtJOvsBw5XXa11kQ\npLeyPh1RydE+WQQN3hMwFp9HwmJF2gzdFvEV5SkjVtB7nIr9m7U6V/TuWGZRCQJ5\neNQjX6f/yb1uTCGgfs0IZNECgYB3yCRcgk+tHfd8dvXPJ09FvAguqisbHgn6oPoo\nUJGdckNdBBVZieaKetQJhPlS50rOfsAnpspVXuQ4FTpJDB9iUjVeuEbF0J8M6Uvj\n6c7jNQKVkhmsIJGrcpkN9+LeiOoNftVVqb55gkYgVD++G0lC+B9cxLyaEQSHnnAV\nV1h2EQKBgCJ6Jb8HXG8g58SF19ksxfDkHoSFJFivwp2Q23tydrofQX/JFKoWvDqP\nUp2C/E0CFr1YcQq4lpJ+QGPUwHOU9xJ4qmSGidBbjFGOLkCLv1Ka20N1rUbcvrnl\nhbxCEp6Ai8WXlJCJ8nTgQR28nAtFsFaPwmvSJAgkADoxf7aYSfVe\n-----END RSA PRIVATE KEY-----"

//// 配置 binlog syncer
//	//cfg := replication.BinlogSyncerConfig{
//	//	ServerID: 100,     // 不能与主库或其他监听者的ServerID冲突
//	//	Flavor:   "mysql", // mysql或mariadb
//	//	Host:     "127.0.0.1",
//	//	Port:     3306,
//	//	User:     "repl",
//	//	Password: "password",
//	//}
//	//
//	//syncer := replication.NewBinlogSyncer(cfg)
//	//
//	//// 可从指定位置启动，也可从最新处启动
//	//streamer, err := syncer.StartSync(mysql.Position{Name: "", Pos: 4})
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//for {
//	//	ev, err := streamer.GetEvent(nil)
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//	// 处理每个binlog事件
//	//	ev.Dump(os.Stdout)
//	//}

//func TestComponent(){
//	dataBase, _ := db.InitDB()
//
//	c := ""
//	list := make([]models.User, 0)
//
//	for i := 0; i < 16; i++ {
//
//		cond := condition.NewConditionBuilder().Build()
//
//		sql := queryexec.NewQueryExecOnMySQL[models.User, uint64](
//			cond,
//			20,
//			nil,
//			dataBase,
//			c,
//		)
//
//		list, c, _, _ = sql.QueryExec(false)
//
//		for _, v := range list {
//			marshal, _ := json.Marshal(v)
//			fmt.Println(string(marshal))
//		}
//
//		fmt.Println("-----------------------------------------------------")
//	}
//
//	for i := 0; i < 16; i++ {
//
//		cond := condition.NewConditionBuilder().Build()
//
//		sql := queryexec.NewQueryExecOnMySQL[models.User, uint64](
//			cond,
//			20,
//			nil,
//			dataBase,
//			c,
//		)
//
//		list, c, _, _ = sql.QueryExec(true)
//
//		for _, v := range list {
//			marshal, _ := json.Marshal(v)
//			fmt.Println(string(marshal))
//		}
//
//		fmt.Println("-----------------------------------------------------")
//	}
//
//	for i := 0; i < 16; i++ {
//
//		cond := condition.NewConditionBuilder().Build()
//
//		sql := queryexec.NewQueryExecOnMySQL[models.User, uint64](
//			cond,
//			20,
//			nil,
//			dataBase,
//			c,
//		)
//
//		list, c, _, _ = sql.QueryExec(false)
//
//		for _, v := range list {
//			marshal, _ := json.Marshal(v)
//			fmt.Println(string(marshal))
//		}
//
//		fmt.Println("-----------------------------------------------------")
//	}
//
//}

func main() {

}
