import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { JelovnikService } from './../../../../services/jelovnik.service';
import { FinansijskaKarticaService } from './../../../../services/finansijskaKartica.service';
import { JelovnikPrikaz } from './../../../../model/JelovnikPrikaz';
import { AuthService } from '../../../../services/auth.service';
@Component({
  selector: 'app-prikaz-jelovnik',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './prikaz-jelovnik.html',
  styleUrls: ['./prikaz-jelovnik.css']
})
export class PrikazJelovnik implements OnInit {
  jelovnici: JelovnikPrikaz[] = [];
  loading = true;
  isAdmin = false;
  userId: string = '';
  role: string = '';

  constructor(
    private jelovnikService: JelovnikService,
    private karticaService: FinansijskaKarticaService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    // 1. Dohvati token i parsiraj payload
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        this.role = payload.role || '';
        this.userId = payload.userId || '';
      } catch (e) {
        console.error('Neuspešno parsiranje tokena', e);
      }
    }

    this.isAdmin = this.authService.isAdmin();

    // 2. Dohvati jelovnike
    this.jelovnikService.getJelovnici().subscribe({
      next: (data) => {
        this.jelovnici = data;
        this.loading = false;
      },
      error: (err) => {
        console.error(err);
        this.loading = false;
      }
    });
  }

  // 3. Funkcija za iskorišćavanje obroka
iskoristiObrok(jelovnikId: string, jeloId: string) {
  if (!this.userId) {
    alert('Niste prijavljeni!');
    return;
  }

  this.karticaService.iskoristiObrok(this.userId, jelovnikId, jeloId).subscribe({
    next: (kartica) => {
      alert('Obrok uspešno iskorišćen!');
      console.log('Nova kartica:', kartica);
    },
    error: (err) => {
      alert('Greška pri iskorišćavanju obroka: ' + err.error);
      console.error(err);
    }
  });
}

}
