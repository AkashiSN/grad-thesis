// convertGoroutine 各Goroutineごとに実行されてURLをicnに変換する
func convertGoroutine(mbuf *util.MovableBuffer, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	var icn strings.Builder     // icnのURLのビルダー
	var etldIcn strings.Builder // ポインタとeTLDとを分割して出力する際のビルダー

	icn.Grow(util.IntermediateMapSize)
	if separate { // ポインタとeTLDとを分割して出力するか
		etldIcn.Grow(util.IntermediateMapSize)
	}

	m := make(map[int]int, util.IntermediateMapSize)
	count := 0

	for _, URL := range mbuf.Buf {
		icnString, segment, isSpecified := generateICN(URL)
		if segment < 0 {
			continue
		}

		if separate && !isSpecified { // ポインタとeTLDとを分割して出力する際のeTLD側
			etldIcn.WriteString(icnString)
		} else { // それ以外はこっち
			icn.WriteString(icnString)
		}

		if v, ok := m[segment]; ok {
			m[segment] = v + 1
		} else {
			m[segment] = 1
		}

		count++
	}

	mbuf.Move() // バッファを開放する

	// mutexをロックする
	mux.Lock()
	defer mux.Unlock()

	icnFileWriter.WriteString(icn.String())
	icnFileWriter.Flush()

	if separate { // ポインタとeTLDとを分割して出力する場合
		etldIcnFileWriter.WriteString(etldIcn.String())
		etldIcnFileWriter.Flush()
	}

	icnCount += count
	for k, v := range m {
		if rv, ok := numberOfSegmentMap[k]; ok {
			numberOfSegmentMap[k] = rv + v
		} else {
			numberOfSegmentMap[k] = v
		}
	}

}

// ConvertICN URLからICNを生成する
func ConvertICN(uniqueURLFilePath string, specificTLDs []string, com, sep bool, bufSize, maxW int) string {
	separate, bufferSize, maxWorkers = sep, bufSize, maxW

	start := time.Now()

	if com { // comの10件以上出現したものをポインタとする
		icnFileName = util.GetFileName(uniqueURLFilePath) + "-icn-com10.txt"
	} else if specificTLDs[0] != "" { // 上位256をポインタとするeTLDが指定されている場合
		icnFileName = util.GetFileName(uniqueURLFilePath) + "-icn-" + strings.Join(specificTLDs, "-") + "256.txt"
	} else { // eTLDのみをポインタとする場合
		icnFileName = util.GetFileName(uniqueURLFilePath) + "-icn.txt"
	}

	graphFile := util.GetFileName(icnFileName) + "-segment.tsv"

	if specificTLDs[0] != "" {
		specifiedRoots = countSpecifiedRoots(uniqueURLFilePath, specificTLDs, com) // 指定したeTLDを条件に沿ったRootを調べる
	}

	log.Printf("Input file : %s", uniqueURLFilePath)
	scanner, inputFile := util.InputFile(uniqueURLFilePath)
	defer inputFile.Close()

	if separate { // ポインタとeTLDとを分割して出力する場合
		etldIcnFileName = util.GetFileName(icnFileName) + "-etld.txt"
		icnFileName = util.GetFileName(icnFileName) + "-pointer.txt"
		etldIcnFileWriter, etldIcnFile = util.OutputFile(etldIcnFileName)
		log.Printf("Output file : %s", etldIcnFileName)
	}
	defer etldIcnFile.Close()

	log.Printf("Output file : %s", icnFileName)
	icnFileWriter, icnFile = util.OutputFile(icnFileName)
	defer icnFile.Close()

	log.Printf("Output segment file : %s", graphFile)

	// tld parser
	var err error
	extract, err = tldextract.New(util.Cache, false)
	if err != nil {
		util.Elog.Print(err)
		os.Exit(1)
	}

	urlCount = util.Parallel(scanner, bufferSize, maxWorkers, convertGoroutine)

	icnFileWriter.Flush()
	if separate { // ポインタとeTLDとを分割して出力する場合
		etldIcnFileWriter.Flush()
	}

	processedURLCount := 0
	segment := []segmentStruct{}
	for s, c := range numberOfSegmentMap {
		segment = append(segment, segmentStruct{
			segment: s,
			count:   c,
		})
		processedURLCount += c
	}
	numberOfSegmentMap = map[int]int{}

	sort.Slice(segment, func(i, j int) bool {
		return segment[i].segment < segment[j].segment
	})

	writer, outFile := util.OutputFile(graphFile)
	for _, v := range segment {
		writer.WriteString(strconv.Itoa(v.segment))
		writer.WriteString("\t")
		writer.WriteString(strconv.Itoa(v.count))
		writer.WriteString("\n")
	}
	writer.Flush()
	outFile.Close()

	os.Remove(util.Cache)

	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of input URLs : %v", humanize.Comma(int64(urlCount)))
	log.Printf("Number of processed URLs : %v", humanize.Comma(int64(processedURLCount)))
	log.Printf("Number of ICNs : %v", humanize.Comma(int64(icnCount)))

	return icnFileName
}
