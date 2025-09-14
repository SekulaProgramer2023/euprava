import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Review {
  id?: string;
  soba_id: string;
  user_id: string;
  rating: number;
  comment?: string;
}

export interface AverageRatingResponse {
  soba_id: string;
  average_rating: number;
}

@Injectable({
  providedIn: 'root'
})
export class ReviewService {

  private apiUrl = 'http://localhost/domovi/reviews'; // promeni ako je drugi port/URL

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
  getReviewsBySoba(sobaId: string): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/${sobaId}`);
  }

  // Dohvatanje prosečne ocene za sobu
  getAverageRating(sobaId: string): Observable<AverageRatingResponse> {
    return this.http.get<AverageRatingResponse>(`${this.apiUrl}/average/${sobaId}`);
  }
}
