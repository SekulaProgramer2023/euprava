import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../model/User'
import { FinansijskaKartica } from '../model/FinansijskaKartica'
import { Dogadjaj } from '../model/dogadjaj.model';

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
  private baseUrl = 'http://localhost:81/menza/users';
   private karticaUrl = 'http://localhost:81/menza/finansijskaKartica/kartice'; 
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
// 🔹 Dodaj jednu ili više alergija
addAlergija(userId: string, alergije: string[]): Observable<User> {
  return this.http.put<User>(`${this.baseUrl}/users/${userId}/alergije`, alergije);
}

// 🔹 Dodaj jedno ili više omiljenih jela
addOmiljenoJelo(userId: string, jela: string[]): Observable<User> {
  return this.http.put<User>(`${this.baseUrl}/users/${userId}/jela`, jela);
}
getEmailFromToken(): string | null {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.email || null;
    } catch (e) {
      return null;
    }
  }

  getUserByEmail(email: string): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/by-email`, { email });
  }
   getKarticaByUser(userId: string): Observable<FinansijskaKartica> {
    return this.http.get<FinansijskaKartica>(`${this.karticaUrl}/user/${userId}`);
  }

  getDogadjaji(): Observable<Dogadjaj[]> {
    return this.http.get<Dogadjaj[]>(`http://localhost/domovi/dogadjaj/dogadjaji`);
  }

  updateStatus(id: string, status: string): Observable<any> {
    return this.http.put(`${this.baseUrl}/dogadjaji/${id}/status`, { status });
  }

}
