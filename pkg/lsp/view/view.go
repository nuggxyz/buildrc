package view

// type view struct {

// 	// keep track of files by document_uri and by basename, a single file may be mapped
// 	// to multiple document_uris, and the same basename may map to multiple files
// 	filesByURI  map[defines.DocumentUri]HCLFile
// 	filesByBase map[string][]HCLFile
// 	fileMu      *sync.RWMutex

// 	openFiles  map[defines.DocumentUri]bool
// 	openFileMu *sync.RWMutex

// 	pbHeaders map[defines.DocumentUri][]string
// }

// func (v *view) GetFile(document_uri defines.DocumentUri) (HCLFile, error) {
// 	if f, ok := v.filesByURI[document_uri]; ok {
// 		return f, nil
// 	}
// 	// no file load try again
// 	err := v.loadHCLFile(document_uri)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if f, ok := v.filesByURI[document_uri]; ok {
// 		return f, nil
// 	}

// 	return nil, fmt.Errorf("%v not found", document_uri)
// }

// // setContent sets the file contents for a file.
// func (v *view) setContent(ctx context.Context, document_uri defines.DocumentUri, data []byte) {

// 	v.fileMu.Lock()
// 	defer v.fileMu.Unlock()

// 	if data == nil {
// 		delete(v.filesByURI, document_uri)
// 		return
// 	}

// 	pf := &hclFile{
// 		File: &file{
// 			document_uri: document_uri,
// 			data:         data,
// 			hash:         hashContent(data),
// 		},
// 	}
// 	pre := v.filesByURI[document_uri]
// 	pf.ref = pre.HCL()
// 	v.filesByURI[document_uri] = pf
// 	// TODO:
// 	//  Control times of parse of ref.
// 	//  Currently it parses every time of file change.
// 	ref, err := parseHCL(document_uri, data)

// 	if err != nil {
// 		return
// 	}
// 	pf.ref = ref
// }

// func (v *view) shutdown(ctx context.Context) error {
// 	// return ViewManagerInstance.RemoveView(ctx, v)
// 	return nil
// }

// func (v *view) didOpen(document_uri defines.DocumentUri, text []byte) {
// 	v.openFileMu.Lock()
// 	v.openFiles[document_uri] = true
// 	v.openFileMu.Unlock()
// 	v.openFile(document_uri, text)
// 	// not like include
// 	v.parseImportHCL(document_uri)
// }

// func (v *view) didOpenPbHeader(document_uri defines.DocumentUri, text string) {
// 	v.pbHeaders[document_uri] = strings.Split(text, "\n")
// }

// func (v *view) GetPbHeaderLine(document_uri defines.DocumentUri, line int) string {
// 	lines, ok := v.pbHeaders[document_uri]
// 	if !ok || len(lines) <= line {
// 		return ""
// 	}

// 	return lines[line]
// }
// func (v *view) didSave(document_uri defines.DocumentUri) {
// 	v.fileMu.Lock()
// 	if file, ok := v.filesByURI[document_uri]; ok {
// 		file.SetSaved(true)
// 	}
// 	v.fileMu.Unlock()
// }

// func (v *view) didClose(document_uri defines.DocumentUri) {
// 	v.openFileMu.Lock()
// 	delete(v.openFiles, document_uri)
// 	v.openFileMu.Unlock()
// }

// func (v *view) isOpen(document_uri defines.DocumentUri) bool {
// 	v.openFileMu.RLock()
// 	defer v.openFileMu.RUnlock()

// 	open, ok := v.openFiles[document_uri]
// 	if !ok {
// 		return false
// 	}
// 	return open
// }

// func (v *view) openFile(document_uri defines.DocumentUri, data []byte) {
// 	v.fileMu.Lock()
// 	defer v.fileMu.Unlock()

// 	pf := &hclFile{
// 		File: &file{
// 			document_uri: document_uri,
// 			data:         data,
// 			hash:         hashContent(data),
// 		},
// 	}

// 	ref, err := parseHCL(document_uri, data)
// 	if err != nil {
// 		return
// 	}
// 	pf.ref = ref
// 	v.filesByURI[document_uri] = pf
// }

// func (v *view) parseSchemaHCL(document_uri defines.DocumentUri) {
// 	ref_file, err := v.GetFile(document_uri)
// 	if err != nil {
// 		logs.Printf("parseImportHCL GetFile err:%v", err)
// 		return
// 	}

// 	ref_file.HCL().Body.Content()

// 	attrs, diag := ref_file.HCL().Body.JustAttributes()
// 	if diag.HasErrors() {
// 		logs.Printf("parseImportHCL JustAttributes err:%v", diag)
// 		return
// 	}
// 	for _, i := range attrs {
// 		import_uri, err := GetDocumentUriFromImportPath(document_uri, i.HCLImport.Filename)
// 		if err != nil {
// 			logs.Printf("parse import err:%v", err)
// 			continue
// 		}
// 		ref_file, err := v.GetFile(import_uri)
// 		if ref_file == nil {
// 			v.loadHCLFile(import_uri)
// 		}
// 	}
// }

// func (v *view) loadHCLFile(document_uri defines.DocumentUri) error {
// 	data, err := os.ReadFile(uri.URI(document_uri).Filename())

// 	if err != nil {
// 		return fmt.Errorf("read file err:%v", err)
// 	}
// 	if !utf8.Valid(data) {
// 		data = toUtf8(data)
// 	}
// 	v.openFile(document_uri, data)
// 	return nil
// }

// func (v *view) mapFile(document_uri defines.DocumentUri, f HCLFile) {
// 	v.fileMu.Lock()

// 	v.filesByURI[document_uri] = f
// 	basename := filepath.Base(uri.URI(document_uri).Filename())
// 	v.filesByBase[basename] = append(v.filesByBase[basename], f)

// 	v.fileMu.Unlock()
// }

// func newView() *view {
// 	return &view{
// 		filesByURI:  make(map[defines.DocumentUri]HCLFile),
// 		filesByBase: make(map[string][]HCLFile),
// 		fileMu:      &sync.RWMutex{},
// 		openFiles:   make(map[defines.DocumentUri]bool),
// 		openFileMu:  &sync.RWMutex{},
// 		pbHeaders:   make(map[defines.DocumentUri][]string),
// 	}
// }

// var ViewManager *view

// func parseHCL(document_uri defines.DocumentUri, data []byte) (*hcl.File, hcl.Diagnostics) {
// 	parser := hclparse.NewParser()

// 	ref, err := parser.ParseHCL(data, string(document_uri))
// 	// if err != nil {
// 	// 	logs.Printf("parseHCL err %v", err)
// 	// }
// 	return ref, err
// }

// func GetDocumentUriFromImportPath(cwd defines.DocumentUri, import_name string) (defines.DocumentUri, error) {
// 	pos := path.Dir(uri.URI(cwd).Filename())
// 	var res defines.DocumentUri
// 	for path.Clean(pos) != "/" {
// 		abs_name := path.Join(pos, import_name)
// 		_, err := os.Stat(abs_name)
// 		if err == nil {
// 			return defines.DocumentUri(uri.New(path.Clean(abs_name))), nil
// 		}
// 		pos = path.Join(pos, "..")
// 	}
// 	return res, fmt.Errorf("import %v not found", import_name)
// }

// func toUtf8(iso8859_1_buf []byte) []byte {
// 	buf := make([]rune, len(iso8859_1_buf))
// 	for i, b := range iso8859_1_buf {
// 		buf[i] = rune(b)
// 	}
// 	return []byte(string(buf))
// }
// func hashContent(content []byte) string {
// 	return fmt.Sprintf("%x", sha1.Sum(content))
// }

// func didOpen(ctx context.Context, params *defines.DidOpenTextDocumentParams) error {
// 	if IsHCLFile(params.TextDocument.Uri) {
// 		document_uri := params.TextDocument.Uri
// 		text := []byte(params.TextDocument.Text)

// 		ViewManager.didOpen(document_uri, text)

// 		return nil
// 	}

// 	if IsPbHeader(params.TextDocument.Uri) {
// 		ViewManager.didOpenPbHeader(params.TextDocument.Uri, params.TextDocument.Text)
// 	}
// 	return nil
// }

// func didChange(ctx context.Context, params *defines.DidChangeTextDocumentParams) error {
// 	if !IsHCLFile(params.TextDocument.Uri) {
// 		return nil
// 	}

// 	if len(params.ContentChanges) < 1 {
// 		return jsonrpc2.NewError(jsonrpc2.InternalError, "no content changes provided")
// 	}

// 	document_uri := params.TextDocument.Uri
// 	text := params.ContentChanges[0].Text

// 	ViewManager.setContent(ctx, document_uri, []byte(text.(string)))
// 	return nil
// }

// func didClose(ctx context.Context, params *defines.DidCloseTextDocumentParams) error {
// 	if !IsHCLFile(params.TextDocument.Uri) {
// 		return nil
// 	}

// 	document_uri := params.TextDocument.Uri

// 	ViewManager.didClose(document_uri)
// 	ViewManager.setContent(ctx, document_uri, nil)

// 	return nil
// }

// func didSave(_ context.Context, params *defines.DidSaveTextDocumentParams) error {
// 	if !IsHCLFile(params.TextDocument.Uri) {
// 		return nil
// 	}

// 	document_uri := defines.DocumentUri(params.TextDocument.Uri)

// 	ViewManager.didSave(document_uri)

// 	return nil
// }

// func onInitialized(ctx context.Context, req *defines.InitializeParams) (err error) {
// 	return nil
// }

// func onDidChangeConfiguration(ctx context.Context, req *defines.DidChangeConfigurationParams) (err error) {
// 	return nil
// }

// func Init(server *lsp.Server) {
// 	ViewManager = newView()

// 	server.OnInitialized(onInitialized)
// 	server.OnDidChangeConfiguration(onDidChangeConfiguration)
// 	server.OnDidOpenTextDocument(didOpen)
// 	server.OnDidChangeTextDocument(didChange)
// 	server.OnDidCloseTextDocument(didClose)
// 	server.OnDidSaveTextDocument(didSave)
// }

// func IsHCLFile(document_uri defines.DocumentUri) bool {
// 	return strings.HasSuffix(string(document_uri), ".hcl")
// }

// func IsPbHeader(document_uri defines.DocumentUri) bool {
// 	return strings.HasSuffix(string(document_uri), ".pb.h")
// }
