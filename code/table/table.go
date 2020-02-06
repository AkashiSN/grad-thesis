// GenerateTable テーブルを作成する
func GenerateTable(inputFilePath, domainFile, tldDir, hashType string) {
	start := time.Now()

	tldScanner, inputFile := util.InputFile(domainFile)
	defer inputFile.Close()

	hashTableFileName := util.GetFileName(inputFilePath) + "-" + hashType + "-hash_table.tsv"
	hashTableWriter, hashTableFile := util.OutputFile(hashTableFileName)
	defer hashTableFile.Close()

	pointerTableFileName := util.GetFileName(inputFilePath) + "-" + hashType + "-pointer_table.tsv"
	pointerTableWriter, pointerTableFile := util.OutputFile(pointerTableFileName)
	defer pointerTableFile.Close()

	log.Printf("Output hash table file : %s", hashTableFileName)
	log.Printf("Output pointer table file : %s", pointerTableFileName)

	icnCount, writeByteSize, hashCount := 0, 0, 0

	pointerTableWriter.WriteString(hashTableFileName) // メタ情報出力
	pointerTableWriter.WriteString("\t")
	pointerTableWriter.WriteString(hashType)
	pointerTableWriter.WriteString("\n")

	for tldScanner.Scan() {
		txt := tldScanner.Text()
		tld := strings.SplitN(txt, "\t", 3)[2]
		rtld := util.ReverseTLD(tld)

		pointerTableWriter.WriteString(tld)
		pointerTableWriter.WriteString("\t")
		pointerTableWriter.WriteString(strconv.Itoa(writeByteSize))
		pointerTableWriter.WriteString("\n")

		icnScanner, icnInputFile := util.InputFile(filepath.Join(tldDir, rtld+".tsv"))
		for icnScanner.Scan() {
			icn := icnScanner.Text()
			icnCount++
			if strings.Split(icn, "\t")[0] != "-" {
				b1, _ := hashTableWriter.WriteString(icn)
				b2, _ := hashTableWriter.WriteString("\n")
				hashCount++
				writeByteSize += b1
				writeByteSize += b2
			}
		}

		hashTableWriter.Flush()
		icnInputFile.Close()
	}

	pointerTableWriter.Flush()
	end := time.Now()

	log.Printf("Total execution time : %v", humanize.SI(float64(end.Sub(start).Seconds()), "s"))
	log.Printf("Number of input ICNs : %v", humanize.Comma(int64(icnCount)))
	log.Printf("Number of generated Tables : %v", humanize.Comma(int64(hashCount)))

}
