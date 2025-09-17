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
  users: string[];
}

interface User {
  id: string;
  name: string;
  surname: string;
  soba?: string;
  selected?: boolean; // koristimo za checkbox
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
  
  addedUsers: User[] = [];    // već dodati korisnici
  allUsers: User[] = [];      // svi korisnici iz baze
  selectedUsers: string[] = []; // ID-evi čekiranih korisnika

  // Modal za kreiranje događaja
  showModal = false;
  newDogadjaj: Dogadjaj = {
    naziv: '',
    opis: '',
    datum_odrzavanja: '',
    datum_zahteva: new Date().toISOString(),
    status: 'na cekanju',
    tema: '',
    users: [],
  };

  // Modal za dodavanje usera
  showUsersModal = false;
  usersWithRooms: User[] = [];
  selectedDogadjajId: string | null = null;

  role: string = '';
  userId: string = '';

  constructor(private router: Router, private http: HttpClient) {}

  ngOnInit() {
    const token = localStorage.getItem('token');
  if (token) {
    const payload = JSON.parse(atob(token.split('.')[1])); 
    this.role = payload.role;
    this.userId = payload.userId;
  }
    this.fetchDogadjaji();
  }



  // ================== DOGADJAJI ==================
  fetchDogadjaji() {
  this.http.get<Dogadjaj[]>('http://localhost/domovi/dogadjaj/dogadjaji')
    .subscribe({
      next: (data) => {
        if (this.role === 'Admin') {
          // Admin vidi sve događaje
          this.dogadjaji = data;
        } else {
          // Student vidi samo događaje gde je dodat
          this.dogadjaji = data.filter(d => d.users.includes(this.userId));
        }
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
      tema: '',
      users: []
    };
    this.showModal = true;
  }

  closeModal() {
    this.showModal = false;
  }

  submitDogadjaj() {
    const payload = {
      ...this.newDogadjaj,
      datum_odrzavanja: new Date(this.newDogadjaj.datum_odrzavanja).toISOString(),
      datum_zahteva: new Date().toISOString()
    };

    this.http.post('http://localhost/domovi/dogadjaj/dogadjaj', payload)
      .subscribe({
        next: (res) => {
          console.log('Događaj kreiran:', res);
          this.closeModal();
          this.fetchDogadjaji();
        },
        error: (err) => {
          console.error('Greška pri kreiranju događaja:', err);
          alert('Greška pri kreiranju događaja.');
        }
      });
  }

  // ================== USERS ==================
  openUsersModal(dogadjajId: string) {
  this.selectedDogadjajId = dogadjajId;

  // Pronađi događaj po ID-u
  const dogadjaj = this.dogadjaji.find(d => d.id === dogadjajId);
  if (!dogadjaj) {
    alert('Događaj nije pronađen!');
    return;
  }

  // Lista već dodanih korisnika (samo ID-evi)
  const addedUserIds = dogadjaj.users || [];

  // Dohvati sve korisnike iz baze
  this.http.get<User[]>('http://localhost/domovi/users/users')
    .subscribe({
      next: (users) => {
        // Filtriramo samo korisnike koji imaju sobu
        const usersWithRoom = users.filter(u => !!u.soba);

        // Razdvojimo već dodate i ostale
        this.addedUsers = usersWithRoom.filter(u => addedUserIds.includes(u.id));
        this.usersWithRooms = usersWithRoom
          .filter(u => !addedUserIds.includes(u.id))
          .map(u => ({ ...u, selected: false }));
        
        this.showUsersModal = true;
      },
      error: (err) => {
        console.error('Greška pri dohvatanju korisnika:', err);
        alert('Nije moguće učitati korisnike.');
      }
    });
}

  closeUsersModal() {
    this.showUsersModal = false;
    this.selectedDogadjajId = null;
    this.usersWithRooms = [];
  }

  submitUsers() {
    if (!this.selectedDogadjajId) return;

    const selectedUsers = this.usersWithRooms
      .filter(u => u.selected)
      .map(u => u.id);

    this.http.post(
      `http://localhost/domovi/dogadjaj/dogadjaj/${this.selectedDogadjajId}/users`,
      { users: selectedUsers }
    ).subscribe({
      next: () => {
        alert('Korisnici uspešno dodati!');
        this.closeUsersModal();
      },
      error: (err) => {
        console.error('Greška pri dodavanju korisnika:', err);
        alert('Greška pri dodavanju korisnika.');
      }
    });
  }

  // ================== PROFIL ==================
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
