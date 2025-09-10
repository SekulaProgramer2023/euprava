import { Component, OnInit } from '@angular/core';
import { RoomService, Soba2 } from '../../../services/room.service';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { UserService } from '../../../services/user.service'
import { User } from '../../../model/user.model';
import { KvarService } from '../../../services/kvar.service';
import { FormsModule } from '@angular/forms';
import { Kvar } from '../../../model/kvar.model';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  userMap: Record<string, { name: string, surname: string }> = {};
  sobe: any[] = [];
  role: string = '';
  userId: string = '';
  dropdownOpen: boolean = false;
  showKvarModal: boolean = false;
  kvarDescription: string = '';
  selectedSobaId: string | null = null;
  kvarMap: Record<string, Kvar[]> = {};



  showErrorModal: boolean = false;   // kontrola modala
  errorMessage: string = '';         // tekst greške

  constructor(private roomService: RoomService, private router: Router, private userService: UserService, private kvarService: KvarService) {}

  ngOnInit(): void {
  const token = localStorage.getItem('token');
  if (token) {
    const payload = JSON.parse(atob(token.split('.')[1])); 
    this.role = payload.role;
    this.userId = payload.userId;
  }

  this.roomService.getSobe().subscribe({
  next: async (data) => {
    this.sobe = data;

    // Inicijalizacija kvarMap za sve sobe
    this.sobe.forEach(soba => {
      this.kvarMap[soba.id] = [];
    });

    // Ako je admin, dohvati korisnike
    if (this.role === 'Admin') {
      const allUserIds = this.sobe.flatMap(soba => soba.users);
      const uniqueUserIds = Array.from(new Set(allUserIds));

      for (const id of uniqueUserIds) {
        await this.userService.getUserById(id).toPromise()
          .then(user => {
            if (user) {
              this.userMap[id] = { 
                name: user?.name ?? 'Nepoznato', 
                surname: user?.surname ?? '' 
              };
            }
          })
          .catch(err => console.error(err));
      }

      // Dohvati kvarove po sobi
      for (const soba of this.sobe) {
        this.kvarService.getKvaroviBySoba(soba.id).subscribe({
          next: (kvarovi) => {
            this.kvarMap[soba.id] = kvarovi;
          },
          error: (err) => console.error('Greška pri dohvatanju kvarova', err)
        });
      }
    }
  },
  error: (err) => this.showError("Greška pri učitavanju soba")
});


}


  useliStudenta(roomId: string): void {
  this.roomService.useliStudenta(roomId, this.userId).subscribe({
    next: (res) => {
      console.log('Uspesno useljen', res);
      // Umesto samo getSobe, updateujemo konkretno sobu u nizu
      const index = this.sobe.findIndex(s => s.id === roomId);
      if (index !== -1) {
        this.sobe[index] = res; // soba sa ažuriranim IsFree i kapacitetom
      }
    },
    error: (err) => {
      console.error('Greška pri useljavanju', err);
      const msg = err.error?.message || err.error || "Došlo je do greške pri useljavanju.";
      this.showError(msg);
    }
  });
}


  showError(message: string) {
    this.errorMessage = message;
    this.showErrorModal = true;
  }

  closeErrorModal() {
    this.showErrorModal = false;
    this.errorMessage = '';
  }

  toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  logout(event?: Event) {
    if(event) event.stopPropagation();
    localStorage.removeItem('token');
    localStorage.removeItem('jwt');
    this.router.navigate(['/domovi/login']);
  }

  goToProfile(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/domovi/profile']);
  }

  // Otvaranje modala
  openKvarModal(sobaId: string) {
    this.selectedSobaId = sobaId;
    this.showKvarModal = true;
  }

  // Zatvaranje modala
  closeKvarModal() {
    this.showKvarModal = false;
    this.kvarDescription = '';
    this.selectedSobaId = null;
  }

  // Slanje prijave
  submitKvar() {
    if (!this.kvarDescription.trim() || !this.selectedSobaId) return;

    this.kvarService.createKvar(this.userId, this.selectedSobaId, this.kvarDescription).subscribe({
      next: () => {
        alert("Kvar uspešno prijavljen!");
        this.closeKvarModal();
      },
      error: (err) => {
        console.error("Greška pri prijavi kvara", err);
        alert("Došlo je do greške pri prijavi kvara.");
      }
    });
  }
}

