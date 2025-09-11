import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';  // ovo je servis, ne u imports
import { RouterModule } from '@angular/router'; // ovo ide u imports
import { CommonModule } from '@angular/common';
import { UserService } from '../../../services/user.service2';
import { FinansijskaKarticaService } from '../../../services/finansijskaKartica.service';
import { FinansijskaKartica } from '../../../model/FinansijskaKartica';
import { HttpErrorResponse } from '@angular/common/http';
import { FormsModule } from '@angular/forms'; // <--- dodaj ovo
@Component({
  selector: 'app-home',
  standalone: true,
 imports: [CommonModule, RouterModule, FormsModule],  // <-- RouterModule, ne Router
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent2 implements OnInit {
  role: string = '';
  userId: string = '';
  dropdownOpen: boolean = false;
  kartica?: FinansijskaKartica;
    depositAmount: number = 0;
  dorucakCount: number = 1;
  rucakCount: number = 1;
  veceraCount: number = 1;
    showHistoryModal: boolean = false;
  constructor(
    private router: Router,   // ovo ostaje ovde
    private karticaService: UserService,
    private karticaService2:FinansijskaKarticaService
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
deposit() {
  if (!this.depositAmount || this.depositAmount <= 0) return;

  this.karticaService2.deposit(this.userId, this.depositAmount).subscribe({
    next: (updated: FinansijskaKartica) => {
      this.kartica = updated;
      this.depositAmount = 0; // reset
    },
    error: (err: any | HttpErrorResponse) => console.error(err)
  });
}

buyDorucak() {
  if (!this.dorucakCount || this.dorucakCount <= 0) return;

  this.karticaService2.buyDorucak(this.userId, this.dorucakCount).subscribe({
    next: (updated: FinansijskaKartica) => {
      this.kartica = updated;
      this.dorucakCount = 1;
    },
    error: (err: any | HttpErrorResponse) => alert(err.error || err.message)
  });
}

buyRucak() {
  if (!this.rucakCount || this.rucakCount <= 0) return;

  this.karticaService2.buyRucak(this.userId, this.rucakCount).subscribe({
    next: (updated: FinansijskaKartica) => {
      this.kartica = updated;
      this.rucakCount = 1;
    },
    error: (err: any | HttpErrorResponse) => alert(err.error || err.message)
  });
}

buyVecera() {
  if (!this.veceraCount || this.veceraCount <= 0) return;

  this.karticaService2.buyVecera(this.userId, this.veceraCount).subscribe({
    next: (updated: FinansijskaKartica) => {
      this.kartica = updated;
      this.veceraCount = 1;
    },
    error: (err: any | HttpErrorResponse) => alert(err.error || err.message)
  });
}


openHistoryModal() {
  this.showHistoryModal = true;
}

closeHistoryModal() {
  this.showHistoryModal = false;
}
}

