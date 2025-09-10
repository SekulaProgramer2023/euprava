import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { FinansijskaKartica } from '../model/FinansijskaKartica';

@Injectable({
  providedIn: 'root'
})
export class FinansijskaKarticaService {
  private baseUrl = 'http://localhost:81/menza/finansijskaKartica';

  constructor(private http: HttpClient) {}

  // Dohvati karticu po userId
  getKarticaByUser(userId: string): Observable<FinansijskaKartica> {
    return this.http.get<FinansijskaKartica>(`${this.baseUrl}/kartice/user/${userId}`);
  }
}
