import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class KvarService {
  private apiUrl = 'http://localhost:80/domovi/sobe'; // prilagodi port

  constructor(private http: HttpClient) {}

  createKvar(userId: string, sobaId: string, description: string): Observable<any> {
    const body = {
      user_id: userId,
      soba_id: sobaId,
      description: description,
      status: false
    };
    return this.http.post(`${this.apiUrl}/prijavi-kvar`, body);
  }

  getKvarovi(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/kvarovi`);
  }

  getKvaroviBySoba(sobaId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/kvarovi/soba/${sobaId}`);
  }

  resolveKvar(kvarId: string) {
  return this.http.put(`${this.apiUrl}/kvarovi/${kvarId}/resolve`, {}); 
}

}
