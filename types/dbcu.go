package types

type DBCU struct {
	Color  uint8 //color for the 1st byte. 8bit color in terminal.
	Color2 uint8 //color for the 2nd byte.
	Rune   rune  //unicode
}
