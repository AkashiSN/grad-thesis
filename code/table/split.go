// SplitByTLD eTLDごとに分割する
func SplitByTLD(icnFileName, hashType string, bufferSize int) string {
	start := time.Now()

	scanner, inputFile := util.InputFile(icnFileName)
	defer inputFile.Close()

	outputDir := util.GetFileName(icnFileName) + "-" + hashType + "-tld/"
	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, os.ModePerm)

	log.Printf("Input file : %s", icnFileName)
	log.Printf("Output dir : %s", outputDir)
	log.Printf("Buffer Size : %v", humanize.Bytes(uint64(bufferSize)))

	originalICNCount, icnCount := 0, 0

	sbuffer := make(map[string]*strings.Builder, 10000)

	for scanner.Scan() {
		icn := scanner.Text()

		hashStr, _ := hash.CalcHash(icn, hashType)

		tld := strings.SplitN(icn, "/", 3)[1] // icn:/tld/

		if sb, ok := sbuffer[tld]; ok {
			sb.WriteString(hashStr)
			sb.WriteString("\t")
			sb.WriteString(icn)
			sb.WriteString("\n")
			icnCount++

			if sb.Cap() > bufferSize { // バッファーサイズを超えたらファイルに書き出す
				writer, outFile := util.OutputFileAppend(outputDir + tld + ".tsv")
				writer.WriteString(sb.String())
				writer.Flush()
				outFile.Close()
				sb.Reset()
			}

		} else {
			sbuffer[tld] = &strings.Builder{}

			sb := sbuffer[tld]
			sb.Grow(util.IntermediateMapSize)
			sb.WriteString(hashStr)
			sb.WriteString("\t")
			sb.WriteString(icn)
			sb.WriteString("\n")

			icnCount++
		}
		originalICNCount++
	}

	inputFile.Close()

	for tld, sb := range sbuffer {
		if sb.Len() != 0 {
			writer, outFile := util.OutputFileAppend(outputDir + tld + ".tsv")
			writer.WriteString(sb.String())
			writer.Flush()
			outFile.Close()
			sb.Reset()
		}
	}

	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of original URLs : %v", humanize.Comma(int64(originalICNCount)))
	log.Printf("Number of processed URLs : %v", humanize.Comma(int64(icnCount)))

	return outputDir
}
