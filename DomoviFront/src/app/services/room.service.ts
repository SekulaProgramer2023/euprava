import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Soba {
  id: string;
  roomNumber: string;
  capacity: number;
  users: string[];
  onBudget: boolean;
  IsFree: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class RoomService {
  private baseUrl = 'http://localhost/domovi/sobe';

  constructor(private http: HttpClient) {}

  getSobe(): Observable<Soba[]> {
    return this.http.get<Soba[]>(`${this.baseUrl}/sobe`);
  }

  useliStudenta(roomId: string, userId: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/useliStudenta`, { roomId, userId });
  }
}
