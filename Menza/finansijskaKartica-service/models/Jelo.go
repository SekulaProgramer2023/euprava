package models

type Nutritijenti struct {
	Masti          float64 `json:"masti"`
	Proteini       float64 `json:"proteini"`
	UgljeniHidrati float64 `json:"ugljeniHidrati"`
}

type Jelo struct {
	JeloID       string       `json:"jeloId"`
	Naziv        string       `json:"naziv"`
	Kategorija   string       `json:"kategorija"`
	TipObroka    string       `json:"tipObroka"`
	Kalorije     float64      `json:"kalorije"`
	Nutritijenti Nutritijenti `json:"nutritijenti"`
}

type Jelovnik struct {
	JelovnikID string `json:"jelovnikId"`
	Datum      string `json:"datum"`
	Dorucak    []Jelo `json:"dorucak"`
	Rucak      []Jelo `json:"rucak"`
	Vecera     []Jelo `json:"vecera"`
	Opis       string `json:"opis"`
}
