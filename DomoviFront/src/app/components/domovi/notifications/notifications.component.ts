import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { NotificationService } from '../../../services/notifications.service';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Dogadjaj } from '../../../model/dogadjaj.model';

@Component({
  selector: 'app-notifications',
  standalone: true,
  imports: [RouterModule, CommonModule, HttpClientModule],
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.css']
})
export class NotificationComponent implements OnInit{
  dropdownOpen: boolean = false;
  notifications: any[] = [];
  role: string = '';
  userId: string = '';

  constructor(private router: Router, private notificationService: NotificationService, private http: HttpClient) {}

  ngOnInit(): void {
    const token = localStorage.getItem('token');
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1]));
      this.role = payload.role;
      this.userId = payload.userId;
    }

    this.loadNotifications();
  }

  loadNotifications(): void {
  const token = localStorage.getItem('token');
  let userId = '';
  let role = '';
  if (token) {
    const payload = JSON.parse(atob(token.split('.')[1]));
    userId = payload.userId;
    role = payload.role;
  }

  // Dohvati sve notifikacije
  this.notificationService.getAllNotifications().subscribe({
    next: (data: any[]) => {
      if (role === 'Admin') {
        // Admin vidi sve notifikacije normalno
        this.notifications = data;
      } else {
        // Student vidi samo svoje relevantne notifikacije
        this.http.get<Dogadjaj[]>('http://localhost/domovi/dogadjaj/dogadjaji')
          .subscribe({
            next: (dogadjaji) => {
              this.notifications = data
                .filter(n => {
                  if (n.dogadjaj_id) {
                    // pronadji događaj
                    const dog = dogadjaji.find(d => d.id === n.dogadjaj_id);
                    if (!dog) return false;
                    // proveri da li je user dodan na događaj
                    return dog.users.includes(userId);
                  }
                  // Ostale notifikacije proveri po user_id
                  return n.user_id === userId;
                })
                .map(n => {
                  if (n.dogadjaj_id) {
                    // Ako je poruka promena statusa događaja, zadrži original
                    if (n.message && n.message.startsWith("Status događaja")) {
                      return n;
                    }
                    // Inače, prikazi generičku poruku za korisnika
                    return { ...n, message: "Dodati ste na dogadjaj" };
                  }
                  return n;
                });
            },
            error: (err) => console.error('Greška pri dohvatanju događaja', err)
          });
      }
    },
    error: (err) => console.error('Greška pri dohvatanju notifikacija', err)
  });
}




  goHome() {
    this.router.navigate(['/domovi/home']);
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
}
