import { Component, OnInit } from '@angular/core';
import { RoomService, Soba2 } from '../../../services/room.service';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { UserService } from '../../../services/user.service'
import { User } from '../../../model/user.model';
import { KvarService } from '../../../services/kvar.service';
import { FormsModule } from '@angular/forms';
import { Kvar } from '../../../model/kvar.model';
import { ReviewService } from '../../../services/review.service';
import { Review } from '../../../model/review.model'; 

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
  selectedKvar: Kvar | null = null;
  showKvarDetailModal: boolean = false;

  showAddReviewModal: boolean = false;

newReview2: { rating1: number; comment: string } = {
  rating1: 1,
  comment: ''
};

averageRatingMap: { [sobaId: string]: number | null } = {};

  reviewsForSoba: Review[] = [];
showReviewsModal = false;

  showErrorModal: boolean = false;   // kontrola modala
  errorMessage: string = '';         // tekst greške

  constructor(private roomService: RoomService, private router: Router, private userService: UserService, private kvarService: KvarService, private reviewService: ReviewService) {}

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
     this.fetchAverageRatings();
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

  goToNotifications(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/domovi/notifications']);
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

  resolveKvar(sobaId: string, kvarId: string) {
  this.kvarService.resolveKvar(kvarId).subscribe({
    next: () => {
      // Ažuriraj status u lokalnom kvarMap da ne mora reload
      const kvarovi = this.kvarMap[sobaId];
      const target = kvarovi.find(k => k.id === kvarId);
      if (target) target.status = true;

      alert("Kvar označen kao rešen");
    },
    error: (err) => {
      console.error("Greška pri označavanju kvara", err);
      alert("Došlo je do greške.");
    }
  });
}

openKvarDetailModal(kvar: Kvar) {
  this.selectedKvar = kvar;
  this.showKvarDetailModal = true;
}

closeKvarDetailModal() {
  this.selectedKvar = null;
  this.showKvarDetailModal = false;
}

resolveSelectedKvar() {
  if (this.selectedKvar) {
    this.resolveKvar(this.selectedKvar.sobaId, this.selectedKvar.id);
    this.closeKvarDetailModal();
  }
}

fetchAverageRatings() {
  this.sobe.forEach(soba => {
    this.reviewService.getAverageRating(soba.id).subscribe({
      next: (res: any) => {
        // res.average_rating je float ili 0
        this.averageRatingMap[soba.id] = res.average_rating > 0 ? res.average_rating : null;
      },
      error: () => {
        this.averageRatingMap[soba.id] = null;
      }
    });
  });
}

openReviewsModal(sobaId: string) {
  this.selectedSobaId = sobaId;
  this.showReviewsModal = true;

  // Dohvati review-e za sobu
  this.reviewService.getReviewsBySoba(sobaId).subscribe({
    next: async (res: Review[]) => {
      this.reviewsForSoba = res;

      // Za svaki review dohvatiti ime i prezime korisnika
      for (const r of this.reviewsForSoba) {
        // Ako već imamo usera u mapi, preskoči
        if (!this.userMap[r.user_id]) {
          try {
            const user = await this.userService.getUserById(r.user_id).toPromise();
            if (user) {
              this.userMap[r.user_id] = { name: user.name ?? 'Nepoznato', surname: user.surname ?? '' };
            } else {
              this.userMap[r.user_id] = { name: 'Nepoznato', surname: '' };
            }
          } catch (err) {
            console.error('Greška pri dohvatanju korisnika za review', err);
            this.userMap[r.user_id] = { name: 'Nepoznato', surname: '' };
          }
        }
      }
    },
    error: () => {
      this.reviewsForSoba = [];
    }
  });
}


closeReviewsModal() {
  this.showReviewsModal = false;
  this.reviewsForSoba = [];
  this.selectedSobaId = null;
}

openAddReviewModal() {
  this.newReview2 = { rating1: 1, comment: '' };
  this.showAddReviewModal = true;
}

closeAddReviewModal() {
  this.showAddReviewModal = false;
}

submitReview() {
  if (!this.selectedSobaId) return;

  const body = {
    soba_id: this.selectedSobaId,
    user_id: this.userId,
    rating: this.newReview2.rating1,
    comment: this.newReview2.comment
  };

  this.reviewService.createReview(body).subscribe({
    next: () => {
      alert("Recenzija uspešno dodata");
      this.closeAddReviewModal();
      this.openReviewsModal(this.selectedSobaId!); // ponovo učitaj review-e
    },
    error: (err) => {
      console.error("Greška pri dodavanju recenzije", err);
      alert("Došlo je do greške pri dodavanju recenzije.");
    }
  });
}

isUserInSelectedRoom(): boolean {
  if (!this.selectedSobaId) return false;

  const soba = this.sobe.find(s => s.id === this.selectedSobaId);
  return soba ? soba.users.includes(this.userId) : false;
}

hasUserReviewed(): boolean {
  if (!this.selectedSobaId) return false;
  return this.reviewsForSoba.some(r => r.user_id === this.userId);
}

}



