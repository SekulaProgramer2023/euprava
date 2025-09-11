export interface IskoriscenoJelo {
  jeloId: string;
  naziv: string;
  tipObroka: 'dorucak' | 'rucak' | 'vecera';
  datum: string; // ISO string iz backend-a
}

export interface FinansijskaKartica {
  id: string;
  userId: string;
  ime: string;
  prezime: string;
  index: string;
  novac: number;
  dorucakCount: number;
  rucakCount: number;
  veceraCount: number;
  iskoriscenaJela: IskoriscenoJelo[]; // novo polje
}
