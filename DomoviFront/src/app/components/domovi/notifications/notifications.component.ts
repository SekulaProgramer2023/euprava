import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { NotificationService } from '../../../services/notifications.service';

@Component({
  selector: 'app-notifications',
  standalone: true,
  imports: [RouterModule, CommonModule],
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.css']
})
export class NotificationComponent implements OnInit{
  dropdownOpen: boolean = false;
  notifications: any[] = [];
  role: string = '';
  userId: string = '';

  constructor(private router: Router, private notificationService: NotificationService) {}

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
  this.notificationService.getAllNotifications().subscribe({
    next: (data) => {
      console.log('Sve notifikacije sa servera:', data);
      if (this.role === 'Admin') {
        this.notifications = data;
      } else {
        this.notifications = data.filter(n =>
          n.user_id?.toString() === this.userId && n.message.toLowerCase().includes('otklonjen')
        );
      }
      console.log('Filtrirane notifikacije:', this.notifications);
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
