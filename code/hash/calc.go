// CalcHash ハッシュを計算する
func CalcHash(icn string, hashT string) (string, bool) {
	sicn := strings.SplitN(icn, "/", 3)
	if len(sicn) < 3 {
		return "-", false
	}
	heads, tails, ok := splitToArrayFillByte(sicn[0] + "/" + sicn[2])
	if !ok {
		return "-", false
	}

	var hash string
	switch hashT {
	case "hash5":
		hash = hash5(heads, tails)
	case "hash6":
		hash = hash6(heads, tails)
	case "hash7":
		hash = hash7(heads, tails)
	}

	return hash, true
}

// calcHashGoroutine GoroutineごとにCalcHashを実行する
func calcHashGoroutine(mbuf *util.MovableBuffer, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	var sbuffer strings.Builder
	sbuffer.Grow(util.IntermediateMapSize)
	count := 0

	for _, icn := range mbuf.Buf {
		hash, ok := CalcHash(icn, hashType)
		if !ok {
			continue
		}
		sbuffer.WriteString(hash)
		sbuffer.WriteString("\t")
		sbuffer.WriteString(icn)
		sbuffer.WriteString("\n")

		count++
	}

	mbuf.Move() // バッファを開放する

	// mutexをロックする
	mux.Lock()
	defer mux.Unlock()

	writer.WriteString(sbuffer.String())
	hashCount += count
	writer.Flush()
}

// Hash ハッシュを計算して出力する
func Hash(inputFilePath, hashT string, bufferSize, maxWorkers int) {
	hashType = hashT

	start := time.Now()

	scanner, inputFile := util.InputFile(inputFilePath)
	defer inputFile.Close()

	outputFileName := util.GetFileName(inputFilePath) + "-" + hashType + ".txt"
	writer, outputFile = util.OutputFile(outputFileName)
	defer outputFile.Close()

	log.Printf("Input file : %s", inputFilePath)
	log.Printf("Output file : %s", outputFileName)
	log.Printf("Buffer Size : %v", humanize.Bytes(uint64(bufferSize)))

	hashCount = 0
	icnCount = util.Parallel(scanner, bufferSize, maxWorkers, calcHashGoroutine)

	writer.Flush()

	end := time.Now()
	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of ICNs : %v", humanize.Comma(int64(icnCount)))
	log.Printf("Number of Hashes : %v", humanize.Comma(int64(hashCount)))
}
