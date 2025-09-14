import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Review {
  id?: string;
  jeloId: string;
  user_id: string;
  rating: number;
  comment?: string;
}

export interface AverageRatingResponse {
  jeloId: string;
  average_rating: number;
}

@Injectable({
  providedIn: 'root'
})
export class ReviewService2 {

  private apiUrl = 'http://localhost:81/menza/review'; // promeni ako je drugi port/URL

  constructor(private http: HttpClient) { }

  // Kreiranje review-a
  createReview(review: Review): Observable<Review> {
    return this.http.post<Review>(`${this.apiUrl}`, review);
  }

  // Dohvatanje svih review-a
  getAllReviews(): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}`);
  }

  // Dohvatanje review-a po sobi
  getReviewsByJelo(jeloId: string): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/${jeloId}`);
  }

  // Dohvatanje prosečne ocene za sobu
  getAverageRating(jeloId: string): Observable<AverageRatingResponse> {
    return this.http.get<AverageRatingResponse>(`${this.apiUrl}/average/${jeloId}`);
  }
}
