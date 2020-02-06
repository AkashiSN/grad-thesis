// SplitByLength 文字列の長さごとに分割する
func SplitByLength(inputFilePath string, freeMem uint64) {
	processFilePath = inputFilePath
	scanner, inputFile := InputFile(processFilePath)
	defer inputFile.Close()

	os.RemoveAll(splitTmpDIr)
	os.MkdirAll(splitTmpDIr, os.ModePerm)

	log.Printf("Input file : %s", processFilePath)
	log.Printf("Output dir : %s", splitTmpDIr)

	bufferSize := freeMem / maxURLLength * 10
	log.Printf("Buffer Size : %v", humanize.Bytes(bufferSize))

	sbuffer := make(map[int]*strings.Builder, maxURLLength)
	start := time.Now()

	for scanner.Scan() {
		url := scanner.Text()
		url = RemoveMeta(url)
		if len(url) == 0 {
			continue
		}

		if sb, ok := sbuffer[len(url)]; ok {
			sb.WriteString(url)
			sb.WriteString("\n")
			processedURLCount++

			if sb.Cap() > int(bufferSize) { // バッファーサイズを超えたらファイルに書き出す
				writer, outFile := OutputFileAppend(splitTmpDIr + strconv.Itoa(len(url)) + ".txt")
				writer.WriteString(sb.String())
				writer.Flush()
				outFile.Close()
				sb.Reset()
			}
		} else {
			sbuffer[len(url)] = &strings.Builder{}

			sb := sbuffer[len(url)]
			sb.Grow(IntermediateMapSize)
			sb.WriteString(url)
			sb.WriteString("\n")

			processedURLCount++
		}
		inputURLCount++
	}

	inputFile.Close()

	for length, sb := range sbuffer {
		if sb.Len() != 0 {
			writer, outFile := OutputFileAppend(splitTmpDIr + strconv.Itoa(length) + ".txt")
			writer.WriteString(sb.String())
			writer.Flush()
			outFile.Close()
			sb.Reset()
		}
	}

	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of input URLs : %v", humanize.Comma(int64(inputURLCount)))
	log.Printf("Number of processed URLs : %v", humanize.Comma(int64(processedURLCount)))

	inputURLCount, processedURLCount = 0, 0
}

// uniqueGoroutine Goroutineごとに実行されて重複を取り除く
func uniqueGoroutine(urlCount, originalURLCount *int, inputFilePath string, writer *bufio.Writer, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	scanner, inFile := InputFile(inputFilePath)
	defer inFile.Close()

	// 1st level uniq per goroutines
	m := make(map[string]bool, IntermediateMapSize)
	count := 0
	for scanner.Scan() {
		url := scanner.Text()
		if len(url) == 0 {
			continue
		}
		count++
		if !m[url] {
			m[url] = true
		}
	}

	// 2nd level uniq with getting mutex
	mux.Lock()
	defer mux.Unlock()

	(*originalURLCount) += count
	for url := range m {
		writer.WriteString(url)
		writer.WriteString("\n")
		(*urlCount)++
	}
	writer.Flush()
}

// Unique 重複をなくす
func Unique(maxWorkers int) {
	files, err := ioutil.ReadDir(splitTmpDIr)
	if err != nil {
		Elog.Print(splitTmpDIr + " is not directory.")
		os.Exit(1)
	}

	outputFilePath := GetFileName(processFilePath) + ".unique"
	writer, outFile := OutputFile(outputFilePath)
	defer outFile.Close()

	log.Printf("Input dir : %s", splitTmpDIr)
	log.Printf("Output file : %s", outputFilePath)

	// sync primitives
	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)

	start := time.Now()

	fileList := []int{}
	for _, file := range files {
		name, _ := strconv.Atoi(GetFileName(file.Name()))
		fileList = append(fileList, name)
	}
	sort.SliceStable(fileList, func(i, j int) bool {
		return fileList[i] < fileList[j]
	})

	for _, file := range fileList {
		inputFilePath := filepath.Join(splitTmpDIr, strconv.Itoa(file)+".txt")
		wg.Add(1)
		go uniqueGoroutine(&processedURLCount, &inputURLCount, inputFilePath, writer, wg, mux)

		// wait if number of goroutines reach max workers for resource limitation
		if runtime.NumGoroutine() >= maxWorkers {
			wg.Wait()
		}
	}

	wg.Wait()
	writer.Flush()

	os.RemoveAll(splitTmpDIr)

	end := time.Now()
	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of input URLs : %v", humanize.Comma(int64(inputURLCount)))
	log.Printf("Number of unique-URLs : %v", humanize.Comma(int64(processedURLCount)))
}
