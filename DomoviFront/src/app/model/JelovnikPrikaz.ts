import { Jelo } from './Jelo';

export interface JelovnikPrikaz {
  jelovnikId?: string;
  datum: string;
  dorucak?: Jelo[];   // niz objekata Jelo
  rucak?: Jelo[];
  vecera?: Jelo[];
  opis?: string;
  
}
