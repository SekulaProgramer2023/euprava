import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Notification {
  _id: string;
  user_id?: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export interface Notification2 {
  _id: string;
  user_id?: string;
  message: string;
  is_read: boolean;
  created_at: string;
  dogadjaj_id: string;
}

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  private baseUrl = 'http://localhost/domovi/notifications';

  constructor(private http: HttpClient) {}

  // Sve notifikacije (za admina)
  getAllNotifications(): Observable<Notification[]> {
    return this.http.get<Notification[]>(`${this.baseUrl}/notifications`);
  }

  // Notifikacije po userId
  getNotificationsByUser(userId: string): Observable<Notification[]> {
    return this.http.get<Notification[]>(`${this.baseUrl}/user/${userId}`);
  }

  // Kreiranje nove notifikacije
  createNotification(notif: Partial<Notification>): Observable<any> {
    return this.http.post(`${this.baseUrl}/notification`, notif);
  }
}
