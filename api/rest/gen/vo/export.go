package vo

type ExportVo struct {
	Title   string   `json:"title"`
	Name    string   `json:"name"`
	Module  string   `json:"module"`
	Methods []string `json:"methods"`
}
