import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';

interface Dogadjaj {
  id?: string;
  naziv: string;
  opis: string;
  datum_odrzavanja: string;
  datum_zahteva: string;
  status: string;
  tema: string;
}

@Component({
  selector: 'app-dogadjaj',
  standalone: true,
  imports: [CommonModule, RouterModule, HttpClientModule, FormsModule],
  templateUrl: './dogadjaj.component.html',
  styleUrls: ['./dogadjaj.component.css']
})
export class DogadjajComponent implements OnInit {
  dropdownOpen = false;
  dogadjaji: Dogadjaj[] = [];
  loading = true;
  error: string | null = null;

  showModal = false;
  newDogadjaj: Dogadjaj = {
    naziv: '',
    opis: '',
    datum_odrzavanja: '',
    datum_zahteva: new Date().toISOString(),
    status: 'na cekanju',
    tema: '',
  };

  constructor(private router: Router, private http: HttpClient) {}

  ngOnInit() {
    this.fetchDogadjaji();
  }

  fetchDogadjaji() {
    this.http.get<Dogadjaj[]>('http://localhost/domovi/dogadjaj/dogadjaji')
      .subscribe({
        next: (data) => {
          this.dogadjaji = data;
          this.loading = false;
        },
        error: (err) => {
          console.error('Greška pri dohvatanju događaja:', err);
          this.error = 'Nije moguće učitati događaje.';
          this.loading = false;
        }
      });
  }

  openModal() {
    this.newDogadjaj = {
      naziv: '',
      opis: '',
      datum_odrzavanja: '',
      datum_zahteva: new Date().toISOString(),
      status: 'na cekanju',
      tema: ''
    };
    this.showModal = true;
  }

  closeModal() {
    this.showModal = false;
  }

  submitDogadjaj() {
  // Pre slanja konvertujemo datum u ISO 8601 format
  const payload = {
    ...this.newDogadjaj,
    datum_odrzavanja: new Date(this.newDogadjaj.datum_odrzavanja).toISOString(),
    datum_zahteva: new Date().toISOString() // sigurnije da backend dobije tačan timestamp
  };

  this.http.post('http://localhost/domovi/dogadjaj/dogadjaj', payload)
    .subscribe({
      next: (res) => {
        console.log('Događaj kreiran:', res);
        this.closeModal();
        this.fetchDogadjaji(); // osveži listu
      },
      error: (err) => {
        console.error('Greška pri kreiranju događaja:', err);
        alert('Greška pri kreiranju događaja.');
      }
    });
}


  toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  logout(event?: Event) {
    if (event) event.stopPropagation();
    localStorage.removeItem('token');
    localStorage.removeItem('jwt');
    this.router.navigate(['/domovi/login']);
  }

  goToProfile(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/domovi/profile']);
  }

  goToNotifications(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/domovi/notifications']);
  }

}
