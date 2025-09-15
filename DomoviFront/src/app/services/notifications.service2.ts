import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
export interface Notification {
  id: string;
  title: string;
  message: string;
  type: string;
  jelovnikID: string;
  jelovnikNaziv?: string;
  jeloID: string;
  jeloNaziv?: string;
  datum_jelovnika:  string;
  is_read: boolean;
  createdAt: string;
}



@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  private baseUrl = 'http://localhost:81/menza/notification';

  constructor(private http: HttpClient) {}

  // Sve notifikacije (za admina)
  getAllNotifications(): Observable<Notification[]> {
    return this.http.get<Notification[]>(`${this.baseUrl}/notification`);
  }

  // Notifikacije po userId
 
  // Kreiranje nove notifikacije
  createNotification(notif: Partial<Notification>): Observable<any> {
    return this.http.post(`${this.baseUrl}/jelo-remaining`, notif);
  }
}
