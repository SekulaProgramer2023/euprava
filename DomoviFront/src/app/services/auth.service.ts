import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor() {}

 getUserRoleFromToken(): string | null {
  const token = localStorage.getItem('token');
  if (!token) return null;

  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.role || null;
  } catch (e) {
    console.error('Greška pri parsiranju tokena', e);
    return null;
  }
}

isAdmin(): boolean {
  return this.getUserRoleFromToken() === 'Admin';
}
}
