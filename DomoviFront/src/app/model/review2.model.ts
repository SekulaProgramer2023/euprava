
export interface Review {
  id?: string;        // opcionalno, MongoDB ID
  jeloId: string;    // ID soe
  user_id: string;    // ID korisnika koji je ostavio review
  rating: number;     // ocena (1-5)
  comment?: string;   // komentar, opcionalno
}
