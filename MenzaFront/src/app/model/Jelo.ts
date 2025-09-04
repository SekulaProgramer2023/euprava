// Enumeracija za kategoriju jela
export enum KategorijaJela {
  Meso = "meso",
  Vegetarijansko = "vegetarijansko",
  Kuvano = "kuvano",
  Desert = "desert",
  Predjelo = "predjelo",
  Salata = "salata"
}

// Enumeracija za tip obroka
export enum TipObroka {
  Dorucak = "dorucak",
  Rucak = "rucak",
  Vecera = "vecera"
}

// Model Jelo
export interface Jelo {
  jeloId?: string;
  naziv: string;
  kategorija: KategorijaJela;
  tipObroka: TipObroka;
  kalorije: number;
  nutritijenti: { [key: string]: number }; // više nije opcionalno
}
