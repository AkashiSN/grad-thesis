// countSpecifiedRoots 指定されたeTLDでの該当するRootを返す
func countSpecifiedRoots(inputFilePath string, specificTLDs []string, com bool) map[string][]string {
	specifiedRoots := make(map[string][]string, len(specificTLDs))
	for _, tld := range specificTLDs {
		scanner, file := util.InputFile(domain.Count(inputFilePath, tld, "", false, false, bufferSize, maxWorkers))
		count := 0
		for scanner.Scan() {
			count++
			txt := strings.Split(scanner.Text(), "\t")
			if com {
				count, _ := strconv.Atoi(txt[0])
				if count < 10 { // 10件以上あるもののみ対象
					break
				}
			} else {
				if count > 256 { // 上位256件をポインタとする
					break
				}
			}

			specifiedRoots[tld] = append(specifiedRoots[tld], txt[2])
		}
		file.Close()
	}

	return specifiedRoots
}

// generateICN ICNを生成する
func generateICN(URL string) (string, int, bool) {
	u, err := url.Parse(URL) // url parse
	if err != nil {
		return "", -1, false
	}

	domain := strings.ToLower(u.Hostname())
	d := extract.Extract(domain) // domain parse
	if d.Flag != 1 {
		return "", -1, false
	}

	var icn strings.Builder
	icn.WriteString("icn:/")

	rtld := util.ReverseTLD(d.Tld)
	icn.WriteString(rtld)

	var isSpecified bool
	if roots, ok := specifiedRoots[d.Tld]; ok && util.Contains(roots, d.Root) { // 指定されているeTLDかつそのRootが条件を満たしているか
		isSpecified = true
		icn.WriteString(".")
	} else {
		isSpecified = false
		icn.WriteString("/")
	}

	icn.WriteString(d.Root) // u.Path is include "/"

	var segment int
	if d.Sub != "" {
		var rsub string
		icn.WriteString("/")
		rsub, segment = util.ReverseSubDomain(d.Sub) // segment数はeTLDとRootを除いたときにスラッシュで区切られている個数
		icn.WriteString(rsub)                        // u.Path is include "/"
	}

	if u.Path != "" {
		path := url.QueryEscape(u.Path)                // マルチバイト文字をエスケープする
		path = strings.ReplaceAll(path, "%2F", "/")    // スラッシュは戻す
		segment += len(strings.Split(u.Path, "/")) - 1 // u.Path is include "/"
		icn.WriteString(path)
	}

	icn.WriteString("\n")

	return icn.String(), segment, isSpecified
}
