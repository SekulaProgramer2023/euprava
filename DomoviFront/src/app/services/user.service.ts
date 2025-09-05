import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../model/user.model';

export interface LoginResponse {
  token: string;
  userId?: string;
}

export interface RegisterData {
  name: string;
  surname: string;
  email: string;
  password: string;
  role: string;
  isActive: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private baseUrl = 'http://localhost:80/domovi/users';

  constructor(private http: HttpClient) {}

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.baseUrl}/login`, { email, password });
  }

  register(data: RegisterData): Observable<any> {
    return this.http.post(`${this.baseUrl}/register`, data);
  }

  getUserById(userId: string): Observable<User> {
  return this.http.get<User>(`${this.baseUrl}/users/${userId}`);
}

getUserIdFromToken(): string | null {
  const token = localStorage.getItem('token');
  if (!token) return null;

  try {
    const payload = JSON.parse(atob(token.split('.')[1])); // dekodira middle deo tokena
    return payload.userId || null;
  } catch (e) {
    console.error('Greška pri parsiranju tokena', e);
    return null;
  }
}

}
