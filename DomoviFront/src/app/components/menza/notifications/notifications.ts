import { Component, OnInit } from '@angular/core';
import { Notification, NotificationService } from '../../../services/notifications.service2';
import { CommonModule } from '@angular/common'; // <--- važ
@Component({
  selector: 'app-notifications',
    imports: [CommonModule],
  templateUrl: './notifications.html',
  styleUrls: ['./notifications.css']
  
})
export class Notifications implements OnInit {

  notifications: Notification[] = [];
  loading = false;
  error = '';

  constructor(private notifService: NotificationService) {}

  ngOnInit(): void {
    this.loadNotifications();
  }

  loadNotifications(): void {
    this.loading = true;
    this.notifService.getAllNotifications().subscribe({
      next: (data) => {
        this.notifications = data.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Greška pri učitavanju notifikacija';
        console.error(err);
        this.loading = false;
      }
    });
  }

  markAsRead(notif: Notification): void {
    notif.is_read = true;
    // Ovde možeš dodati poziv prema backend-u da update-uješ is_read
  }

  formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleString(); // lokalno lepo formatiranje
  }
}
