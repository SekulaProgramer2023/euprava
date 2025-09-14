
export interface Review {
  id?: string;        // opcionalno, MongoDB ID
  soba_id: string;    // ID sobe
  user_id: string;    // ID korisnika koji je ostavio review
  rating: number;     // ocena (1-5)
  comment?: string;   // komentar, opcionalno
}
