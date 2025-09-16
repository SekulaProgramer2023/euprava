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
// Uplata novca
deposit(userId: string, novac: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/${userId}/deposit`,
    { novac }
  );
}

// Kupovina doručka
buyDorucak(userId: string, count: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/${userId}/buy/dorucak`,
    { count }
  );
}

// Kupovina ručka
buyRucak(userId: string, count: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/${userId}/buy/rucak`,
    { count }
  );
}

// Kupovina večere
buyVecera(userId: string, count: number): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/${userId}/buy/vecera`,
    { count }
  );
}
// Iskoristi obrok (dorucak, rucak ili vecera, backend sam odlučuje po tipu)
iskoristiObrok(userId: string, jelovnikId: string, jeloId: string): Observable<FinansijskaKartica> {
  return this.http.post<FinansijskaKartica>(
    `${this.baseUrl}/kartice/iskoristi/${userId}/${jelovnikId}/${jeloId}`,
    {}
  );
}
// Statistika o iskorišćenim obrocima
getStatistika(): Observable<any> {
  return this.http.get<any>(`${this.baseUrl}/kartice/statistika`);
}


}
