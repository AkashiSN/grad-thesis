// countHash メインプロセス
func countHash(mbuf *util.MovableBuffer, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	m := make(map[string]int, util.IntermediateMapSize)
	for _, k := range mbuf.Buf {
		key := strings.SplitN(k, "\t", 2)[0]
		if key == "-" {
			continue
		}

		if v, ok := m[key]; ok {
			m[key] = v + 1
		} else {
			m[key] = 1
		}
	}

	mbuf.Move() // バッファを開放する

	// mutexをロックする
	mux.Lock()
	defer mux.Unlock()

	for k, v := range m {
		if rv, ok := result[k]; ok {
			result[k] = rv + v
		} else {
			result[k] = v
		}
	}
}

// Count ハッシュの衝突数を数える
func Count(inputFilePath string, bufferSize, maxWorkers int) {
	start := time.Now()

	scanner, inputFile := util.InputFile(inputFilePath)
	defer inputFile.Close()

	mode := false
	outputFileName := ""
	if strings.Contains(inputFilePath, "-count") {
		outputFileName = util.GetFileName(inputFilePath) + "-ccdf.tsv"
		mode = true
	} else {
		outputFileName = util.GetFileName(inputFilePath) + "-count.tsv"

	}

	writer, outputFile := util.OutputFile(outputFileName)
	defer outputFile.Close()

	log.Printf("Input file : %s", inputFilePath)
	log.Printf("Output file : %s", outputFileName)
	log.Printf("Buffer Size : %v", humanize.Bytes(uint64(bufferSize)))

	// resultMap
	result = make(map[string]int, bufferSize)

	icnCount = util.Parallel(scanner, bufferSize, maxWorkers, countHash)

	if mode {
		hash := []hashCountCountStruct{}

		for h, c := range result {
			h, _ := strconv.Atoi(h)
			hash = append(hash, hashCountCountStruct{
				count:     c,
				hashCount: h,
			})
		}

		result = nil

		sort.Slice(hash, func(i, j int) bool {
			return hash[i].hashCount < hash[j].hashCount
		})

		hashCount = len(hash)

		cdf := 0
		for _, h := range hash {
			cdf += h.count
			writer.WriteString(strconv.Itoa(h.hashCount))
			writer.WriteString("\t")
			writer.WriteString(strconv.Itoa(h.count))
			writer.WriteString("\t")
			writer.WriteString(strconv.Itoa(cdf))
			writer.WriteString("\t")
			writer.WriteString(strconv.FormatFloat(float64((icnCount-cdf))/float64(icnCount), 'E', 6, 64))
			writer.WriteString("\n")
		}
	} else {
		hash := []hashCountStruct{}

		for h, c := range result {
			hash = append(hash, hashCountStruct{
				count: c,
				hash:  h,
			})
		}

		result = nil

		sort.Slice(hash, func(i, j int) bool {
			return hash[i].count > hash[j].count
		})

		hashCount = len(hash)

		for _, h := range hash {
			writer.WriteString(strconv.Itoa(h.count))
			writer.WriteString("\t")
			writer.WriteString(h.hash)
			writer.WriteString("\n")
		}
	}

	writer.Flush()
	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of ICNs : %v", humanize.Comma(int64(icnCount)))
	log.Printf("Number of Hashes : %v", humanize.Comma(int64(hashCount)))
}
