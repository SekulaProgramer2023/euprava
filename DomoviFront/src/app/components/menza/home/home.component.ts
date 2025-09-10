import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';  // ovo je servis, ne u imports
import { RouterModule } from '@angular/router'; // ovo ide u imports
import { CommonModule } from '@angular/common';
import { UserService } from '../../../services/user.service2';
import { FinansijskaKartica } from '../../../model/FinansijskaKartica';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule],  // <-- RouterModule, ne Router
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent2 implements OnInit {
  role: string = '';
  userId: string = '';
  dropdownOpen: boolean = false;
  kartica?: FinansijskaKartica;

  constructor(
    private router: Router,   // ovo ostaje ovde
    private karticaService: UserService
  ) {}

  ngOnInit(): void {
    const token = localStorage.getItem('token');
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1]));
      this.role = payload.role;
      this.userId = payload.userId;

      // Dohvati karticu za korisnika
      this.karticaService.getKarticaByUser(this.userId).subscribe({
        next: (data) => this.kartica = data,
        error: (err) => console.error('Greška pri dohvaćanju kartice:', err)
      });
    }
  }

  toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  logout(event?: Event) {
    if (event) event.stopPropagation();
    localStorage.removeItem('token');
    this.router.navigate(['/menza/login']);
  }

  goToProfile(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/menza/profile']);
  }
}
