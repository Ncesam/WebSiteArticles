package types

type Article struct {
	Title string
	Blocks []Block
}
  
type Block struct {
	Type string
	Data BlockData
}
type BlockData struct {
	Text string
	Lang string
	Items []string
	Type string
}
  