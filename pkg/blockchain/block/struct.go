package block

type Properties map[string]string

type File struct{
	Properties Properties
	RowData []byte	
}

type Block struct{
	Properties Properties
	NumberOfFiles uint64
	Files []File
}
