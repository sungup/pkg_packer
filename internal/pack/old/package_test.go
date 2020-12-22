package old

/*
func createTestScripts(t *testing.T) ([]string, string) {
	scripts := test.RandStrings(t)

	body := make([]string, 0)

	for _, item := range scripts {
		body = append(body, item+";")
	}

	return scripts, strings.Join(body, "\n")
}

func TestLoadPkgInfo(t *testing.T) {
	a := assert.New(t)
	projectHome := test.GetProjectFileHome()

	/////////////////////// Loading test yaml 1 (without source file option)
	// 0. Prepare expect data structure
	expect, err := test.LoadTestYAMLData(test.ExpectFile1)

	a.NoErrorf(err, "reference set loading fail - expect file 1")

	// 1. LoadPkgInfo test
	expectPath, _ := test.GetTestFilePath(test.ExpectFile1)
	tested, err := LoadPkgInfo(expectPath, projectHome)

	a.NoErrorf(err, "unexpected load error")
	a.Equal(projectHome, tested.srcHome)

	// 2. PackageMeta Test
	compareMetaData(t, expect.Meta, &tested.Meta)

	/////////////////////// Loading test yaml 2 (with source file option)
	// 0. Prepare expect data structure
	expect, err = test.LoadTestYAMLData(test.ExpectFile2)

	a.NoErrorf(err, "reference set loading fail - expect file 2")

	// 1. LoadPkgInfo test
	expectPath, _ = test.GetTestFilePath(test.ExpectFile2)
	tested, err = LoadPkgInfo(expectPath, projectHome)

	a.NoErrorf(err, "unexpected load error")
	a.Equal(projectHome, tested.srcHome)

	// 2. Check file paths
	for _, fileList := range tested.Files {
		for _, file := range fileList {
			if file.Src == "" {
				continue
			}

			expectSrcPath := file.Body

			if !strings.HasPrefix(expectSrcPath, "/") {
				expectSrcPath = path.Join(tested.srcHome, expectSrcPath)
			}

			a.Equal(expectSrcPath, file.Src)
		}
	}
}

func TestPackage_PreInScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PreIn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPreIn(item)
	}

	tested := pkg.PreInScript()

	a.Equal(expect, tested)
}

func TestPackage_PostInScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PostIn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPostIn(item)
	}

	tested := pkg.PostInScript()

	a.Equal(expect, tested)
}

func TestPackage_PreUnScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PreUn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPreUn(item)
	}

	tested := pkg.PreUnScript()

	a.Equal(expect, tested)
}

func TestPackage_PostUnScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PostUn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPostUn(item)
	}

	tested := pkg.PostUnScript()

	a.Equal(expect, tested)
}

func TestPackage_AddFile(t *testing.T) {
	var err error

	a := assert.New(t)
	pkg := NewPackage(PackageMeta{}, test.GetProjectFileHome())
	fileType := "generic"

	////////////
	// Index 0. Append generic file from home directory
	err = pkg.AddFile(
		fileType,
		file{Dest: "/tmp/test1.yml", Src: "test/test.yml"},
	)
	a.NoError(err)

	a.Equal(path.Join(pkg.srcHome, "test/test.yml"), pkg.Files[fileType][0].Src)

	////////////
	// Index 1. Append generic file from absolute path
	err = pkg.AddFile(
		fileType,
		file{Dest: "/tmp/test2.yml", Src: "/tmp/temp.yml"},
	)
	a.NoError(err)

	a.Equal("/tmp/temp.yml", pkg.Files[fileType][1].Src)

	////////////
	// Index 2. Append generic file from empty source with body
	err = pkg.AddFile(
		fileType,
		file{Dest: "/tmp/test3.yml", Body: `test: test data`},
	)
	a.NoError(err)

	a.Equal("", pkg.Files[fileType][2].Src)
}
*/
