import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { FinansijskaKartica } from '../model/FinansijskaKartica';
import { map } from 'rxjs/operators';
@Injectable({
  providedIn: 'root'
})
export class FinansijskaKarticaService {
  private baseUrl = 'http://localhost:81/menza/finansijskaKartica';

  constructor(private http: HttpClient) {}

  // Dohvati karticu po userId
getKarticaByUser(userId: string): Observable<FinansijskaKartica> {
  return this.http.get<FinansijskaKartica>(`${this.baseUrl}/kartice/user/${userId}`)
    .pipe(
      map(kartica => ({
        ...kartica,
        novac: +kartica.novac,
        dorucakCount: +kartica.dorucakCount,
        rucakCount: +kartica.rucakCount,
        veceraCount: +kartica.veceraCount
      }))
    );
}

getKarticaByEmail(email: string): Observable<FinansijskaKartica> {
  return this.http.get<FinansijskaKartica>(`${this.baseUrl}/kartice/userE/by-email`, {
    params: { email }   // email kao query parametar
  }).pipe(
    map(kartica => ({
      ...kartica,
      novac: +kartica.novac,
      dorucakCount: +kartica.dorucakCount,
      rucakCount: +kartica.rucakCount,
      veceraCount: +kartica.veceraCount
    }))
  );
}

// Uplata novca
// Uplata novca
deposit(email: string, novac: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/deposit`,
    { novac },
    { params: { email } }   // email kao query param
  );
}

// Kupovina doručka
buyDorucak(email: string, count: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/buy/dorucak`,
    { count },
    { params: { email } }
  );
}

// Kupovina ručka
buyRucak(email: string, count: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/buy/rucak`,
    { count },
    { params: { email } }
  );
}

// Kupovina večere
buyVecera(email: string, count: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/buy/vecera`,
    { count },
    { params: { email } }
  );
}

// Iskoristi obrok (dorucak, rucak ili vecera, backend sam odlučuje po tipu)
iskoristiObrok(email: string, jelovnikId: string, jeloId: string): Observable<FinansijskaKartica> {
  const params = {
    email,
    jelovnikId,
    jeloId
  };

  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/iskoristi`,
    {},   // body je prazan
    { params }  // šalješ query parametre
  );
}

// Statistika o iskorišćenim obrocima
getStatistika(): Observable<any> {
  return this.http.get<any>(`${this.baseUrl}/kartice/statistika`);
}


}
