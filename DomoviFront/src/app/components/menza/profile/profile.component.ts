import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { UserService } from '../../../services/user.service2';
import { JeloService } from '../../../services/jela.service';
import { User } from '../../../model/User';
import { Jelo } from '../../../model/Jelo';
import { Dogadjaj } from '../../../model/dogadjaj.model';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent2 implements OnInit {
  dropdownOpen: boolean = false;
  user: User | null = null;

  dogadjaji: Dogadjaj[] = [];

  // inputi za forme
  newAlergija: string = '';
  newOmiljenoJelo: string = '';

  // lista svih jela za dropdown
  allJela: Jelo[] = [];

  constructor(
    private userService: UserService,
    private jeloService: JeloService,
    private router: Router
  ) {}

  ngOnInit(): void {
    const email = this.userService.getEmailFromToken();
    if (email) {
      this.userService.getUserByEmail(email).subscribe({
        next: (res) => {
          this.user = new User(
            res.id || '',
            res.password || '',
            res.role || '',
            res.name || '',
            res.surname || '',
            res.email || '',
            res.isActive || false,
            res.alergije || [],
            res.omiljenaJela || []
          );
        },
        error: (err) => console.error('Greška pri dohvatanju korisnika', err)
      });

      this.loadDogadjaji();
    }

    // Dohvati sva jela za dropdown
    this.jeloService.getJela().subscribe({
      next: (jela) => (this.allJela = jela),
      error: (err) => console.error('Greška pri dohvatanju jela', err)
    });
  }

  // 🔹 Navbar metode
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

  goHome() {
    this.router.navigate(['/menza/home']);
  }

  // 🔹 Dodavanje alergije
  addAlergija() {
    if (!this.newAlergija || !this.user) return;

    const updatedAlergije = [...(this.user.alergije || []), this.newAlergija];

    this.userService.addAlergija(this.user.id, updatedAlergije).subscribe({
      next: (updatedUser) => {
        if (!this.user) return;
        this.user.alergije = [...(updatedUser.alergije || [])]; // update UI
        this.newAlergija = ''; 
        window.location.reload();// očisti input
      },
      error: (err) => console.error('Greška pri dodavanju alergije', err)
    });
  }

  // 🔹 Dodavanje omiljenog jela iz dropdown
  addOmiljenoJelo() {
    if (!this.newOmiljenoJelo || !this.user) return;

    const updatedJela = [...(this.user.omiljenaJela || []), this.newOmiljenoJelo];

    this.userService.addOmiljenoJelo(this.user.id, updatedJela).subscribe({
      next: (updatedUser) => {
        if (!this.user) return;
        this.user.omiljenaJela = [...(updatedUser.omiljenaJela || [])]; // update UI
        this.newOmiljenoJelo = '';
        window.location.reload(); // očisti dropdown izbor
      },
      error: (err) => console.error('Greška pri dodavanju omiljenog jela', err)
    });
  }
  loadDogadjaji() {
    this.userService.getDogadjaji().subscribe({
      next: (res) => {
        this.dogadjaji = res.filter(d => d.status === 'na čekanju');
      },
      error: (err) => console.error('Greška pri dohvatanju događaja', err)
    });
  }

  updateStatus(id: string, status: string) {
    this.userService.updateStatus(id, status).subscribe({
      next: () => {
        this.loadDogadjaji(); // ponovo učitaj listu
      },
      error: (err) => console.error('Greška pri ažuriranju statusa', err)
    });
  }
}
