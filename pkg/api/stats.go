package api

type ReturnedRaidStats struct {
	Raids   map[string]map[string]int64 `json:"raids"`
	Classes map[string]int64            `json:"classes"`
	Players int64                       `json:"players"`
}

//func (s *Server) statsHandler(c *gin.Context) {
//	var wg sync.WaitGroup
//	var (
//		rows    []sql.GetRaidStatsRow
//		count   []sql.CountClassesRow
//		players int64
//	)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	wg.Add(3)
//
//	go func() {
//		defer wg.Done()
//		var err error
//		rows, err = s.conn.Queries.GetRaidStats(ctx)
//		if err != nil && !errors.Is(err, context.Canceled) {
//			cancel()
//			log.Println(errors.Wrap(err, "getting raid stats"))
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		var err error
//		players, err = s.conn.Queries.GetUniqueUploaders(ctx)
//		if err != nil && !errors.Is(err, context.Canceled) {
//			cancel()
//			log.Println(errors.Wrap(err, "counting unique uploaders"))
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		var err error
//		count, err = s.conn.Queries.CountClasses(ctx)
//		if err != nil && !errors.Is(err, context.Canceled) {
//			cancel()
//			log.Println(errors.Wrap(err, "counting classes"))
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//	}()
//
//	wg.Wait()
//
//	raids := make(map[string]map[string]int64)
//	for _, row := range rows {
//		if _, ok := raids[row.Difficulty]; !ok {
//			raids[row.Difficulty] = make(map[string]int64)
//		}
//
//		name := row.Boss
//		raid, ok := process.RaidLookup[row.Boss]
//		if ok {
//			name = raid[0]
//		}
//
//		raids[row.Difficulty][name] += row.Count
//	}
//
//	classes := make(map[string]int64)
//	for _, entry := range count {
//		classes[entry.Class] = entry.Count
//	}
//
//	c.JSON(http.StatusOK, ReturnedRaidStats{
//		Raids:   raids,
//		Classes: classes,
//		Players: players,
//	})
//}
