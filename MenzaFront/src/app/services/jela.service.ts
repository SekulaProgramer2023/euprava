import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Jelo } from '../model/Jelo';   // ovde importuješ iz foldera model

@Injectable({
  providedIn: 'root'
})
export class JeloService {
  private baseUrl = 'http://localhost:81/menza/jelovnik/jela';

  constructor(private http: HttpClient) {}

  getJela(): Observable<Jelo[]> {
    return this.http.get<Jelo[]>(this.baseUrl);
  }
  getJelaByTip(tip: string): Observable<Jelo[]> {
  return this.http.get<Jelo[]>(`${this.baseUrl}/tip?tip=${tip}`);
}
  createJelo(jelo: Jelo): Observable<Jelo> {
    return this.http.post<Jelo>(this.baseUrl, jelo);
  }


}
