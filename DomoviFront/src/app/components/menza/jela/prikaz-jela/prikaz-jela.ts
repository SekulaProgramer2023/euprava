import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { JeloService } from '../../../../services/jela.service';
import { AuthService } from '../../../../services/auth.service';
import { ReviewService2 } from '../../../../services/review.service2';
import { Jelo } from '../../../../model/Jelo';
import { Review } from '../../../../model/review2.model';
import { UserService } from '../../../../services/user.service2';
import { ChangeDetectorRef } from '@angular/core';
@Component({
  selector: 'app-prikaz-jela',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './prikaz-jela.html',
  styleUrls: ['./prikaz-jela.css']
})
export class PrikazJela implements OnInit {
  jela: Jelo[] = [];
  loading = true;
  tipPretrage: string = '';
  isAdmin = false;
  newRating: number = 0;
  newComment: string = '';  

  averageRatingMap: { [jeloId: string]: number | null } = {};
  reviewsForJelo: Review[] = [];
  showReviewsModal = false;
  selectedJeloId: string | null = null;

  userMap: Record<string, { name: string; surname: string }> = {};

  constructor(
    private jeloService: JeloService,
    private authService: AuthService,
    private reviewService: ReviewService2,
    private userService: UserService,
       private cd :ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.isAdmin = this.authService.isAdmin();
    this.loadAllJela();
  }

  loadAllJela(): void {
    this.loading = true;
    this.jeloService.getJela().subscribe({
      next: (data) => {
        this.jela = data;
        this.loading = false;
        this.fetchAverageRatings();
      },
      error: (err) => {
        console.error(err);
        this.loading = false;
      }
    });
  }

  searchByTip(): void {
    this.loading = true;
    const source$ = this.tipPretrage
      ? this.jeloService.getJelaByTip(this.tipPretrage)
      : this.jeloService.getJela();

    source$.subscribe({
      next: (data) => {
        this.jela = data;
        this.loading = false;
        this.fetchAverageRatings();
      },
      error: (err) => {
        console.error(err);
        this.loading = false;
      }
    });
  }
fetchAverageRatings() {
  this.jela.forEach(jelo => {
    const id = jelo.jeloId;
    if (!id) return; // preskoči ako jeloId nije definisan

    this.reviewService.getAverageRating(id).subscribe({
      next: (res: any) => {
        this.averageRatingMap[id] = res.average_rating > 0 ? res.average_rating : null;
      },
      error: () => {
        this.averageRatingMap[id] = null;
      }
    });
  });
}


  openReviewsModal(jeloId: string) {
    this.selectedJeloId = jeloId;
    this.showReviewsModal = true;

    this.reviewService.getReviewsByJelo(jeloId).subscribe({
      next: async (res: Review[]) => {
        this.reviewsForJelo = res;

        for (const r of this.reviewsForJelo) {
          if (!this.userMap[r.user_id]) {
            try {
              const user = await this.userService.getUserById(r.user_id).toPromise();
              this.userMap[r.user_id] = user
                ? { name: user.name ?? 'Nepoznato', surname: user.surname ?? '' }
                : { name: 'Nepoznato', surname: '' };
            } catch {
              this.userMap[r.user_id] = { name: 'Nepoznato', surname: '' };
            }
          }
        }
      },
      error: () => {
        this.reviewsForJelo = [];
      }
    });
  }

  closeReviewsModal() {
    this.showReviewsModal = false;
    this.reviewsForJelo = [];
    this.selectedJeloId = null;
  }

addReview() {
  if (!this.selectedJeloId) return;

  const token = localStorage.getItem('token');
  if (!token) return;

  const payload = JSON.parse(atob(token.split('.')[1]));
  const currentUserId = payload.userId;

  // VALIDACIJA
  if (this.newRating < 1 || this.newRating > 5) {
    alert('Ocena mora biti između 1 i 5!');
    return;
  }

  const review: Review = {
    jeloId: this.selectedJeloId,
    user_id: currentUserId,
    rating: this.newRating,
    comment: this.newComment
  };

  this.reviewService.createReview(review).subscribe({
    next: (res) => {
      // Dodaj review u lokalnu listu (za modal ili detaljni prikaz)
      this.reviewsForJelo = [...this.reviewsForJelo, res];

      // Reset inputa
      this.newRating = 0;
      this.newComment = '';

      // Zatvori modal
      this.showReviewsModal = false;

      // Osveži Angular prikaz odmah
      this.cd.detectChanges();

      // Ažuriraj prosečnu ocenu u glavnoj listi
      this.fetchAverageRatings();
    },
    error: (err) => {
      console.error('Greška pri dodavanju review-a', err);
      alert('Došlo je do greške pri dodavanju review-a.');
    }
  });
}
}



