package model

type TicketRequest struct {
	Tujuan  string `json:"tujuan"`
	Tanggal string `json:"tanggal"`
	Waktu   string `json:"waktu"`
}
