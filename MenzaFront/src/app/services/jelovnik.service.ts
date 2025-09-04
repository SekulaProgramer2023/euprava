import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { JelovnikPrikaz } from '../model/JelovnikPrikaz';

@Injectable({
  providedIn: 'root'
})
export class JelovnikService {
  private apiUrl = 'http://localhost:81/menza/jelovnik/jelovnici-sa-jelima';

  constructor(private http: HttpClient) {}

  // Dohvata jelovnike sa imenima jela
  getJelovnici(): Observable<JelovnikPrikaz[]> {
    return this.http.get<JelovnikPrikaz[]>(this.apiUrl);
  }
}
