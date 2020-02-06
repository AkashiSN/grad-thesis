// CreatePointer 行数に対応するポインターを作る
func CreatePointer(inputFilePath string) {
	processFilePath = inputFilePath
	scanner, inputFile := InputFile(processFilePath)
	defer inputFile.Close()

	pointerFileName = GetFileName(processFilePath) + "-pointer.bin"
	pointerFileWriter, pointerFile := OutputFile(pointerFileName)

	log.Printf("Input file : %s", processFilePath)
	log.Printf("Output pointer table file : %s", pointerFileName)

	start := time.Now()

	var readByte uint64
	for scanner.Scan() {
		buf := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(buf, readByte)
		pointerFileWriter.Write(buf)
		readByte += uint64(len(scanner.Bytes()))
		readByte++ // for "\n"
		inputURLCount++
	}

	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, readByte)
	pointerFileWriter.Write(buf)

	pointerFileWriter.Flush()
	pointerFile.Close()
	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of input URLs : %v", humanize.Comma(int64(inputURLCount)))
}

// Random ファイルをランダムに分割する
func Random() {
	scanner, inputFile := InputFile(pointerFileName)
	log.Printf("Input pointer file : %s", pointerFileName)
	log.Printf("Number of input URLs : %v", humanize.Comma(int64(inputURLCount)))

	start := time.Now()
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		return binary.MaxVarintLen64, data[0:binary.MaxVarintLen64], nil
	}
	scanner.Split(split)

	pointerList := make([]uint64, inputURLCount+1)
	i := 0
	var pointer uint64
	for scanner.Scan() {
		bytes := scanner.Bytes()
		pointer = binary.LittleEndian.Uint64(bytes)
		pointerList[i] = pointer
		i++
	}

	inputFile.Close()
	randomFileCount, processedURLCount := 0, 0

	outputDir := GetFileName(processFilePath) + "-random/"
	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, os.ModePerm)

	log.Printf("Output dir : %s", outputDir)
	log.Printf("Split Size : %v", humanize.Bytes(uint64(splitSize)))
	writer, outFile := OutputFile(outputDir + strconv.Itoa(randomFileCount) + ".txt")

	fp, err := os.Open(processFilePath)
	if err != nil {
		panic(err)
	}
	usedNum := make(map[int]bool, inputURLCount)

	num := 0
	var writeByteSize int64

	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rng := rand.New(mt19937.New()) // Mersenne twister
	rng.Seed(seed.Int64())

	for ; ; num = rng.Intn(inputURLCount) {
		if usedNum[num] {
			continue
		}
		usedNum[num] = true

		startPointer := pointerList[num]
		length := pointerList[num+1] - startPointer // include "\n"

		fp.Seek(int64(startPointer), 0) // seek to pointer
		buf := make([]byte, length)
		_, err = io.ReadFull(fp, buf)
		if err != nil {
			panic(err)
		}

		processedURLCount++
		writer.Write(buf)
		writeByteSize += int64(length)

		if writeByteSize >= int64(splitSize) {
			writer.Flush()
			outFile.Close()
			randomFileCount++
			if randomFileCount > 30 {
				break
			}
			writer, outFile = OutputFile(outputDir + strconv.Itoa(randomFileCount) + ".txt")
			writeByteSize = 0
		}
	}

	writer.Flush()
	outFile.Close()
	fp.Close()
	os.RemoveAll(pointerFileName)
	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of processed URLs : %v", humanize.Comma(int64(processedURLCount)))
}
