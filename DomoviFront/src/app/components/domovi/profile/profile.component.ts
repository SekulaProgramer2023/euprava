import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
import { UserService } from '../../../services/user.service';
import { User } from '../../../model/user.model';
import { RoomService } from '../../../services/room.service';
import { Soba } from '../../../model/soba.model'

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  dropdownOpen: boolean = false;
  user: User | null = null;
  roomNumber: string | null = null;   // broj sobe

  constructor(
    private userService: UserService, 
    private roomService: RoomService,
    private router: Router
  ) {}

  ngOnInit(): void {
  const email = this.userService.getEmailFromToken();
  if (email) {
    this.userService.getUserByEmail(email).subscribe({
      next: (res) => {
        this.user = res;
        console.log("User data:", this.user);

        if (this.user?.soba) {
          this.roomService.getSobaById(this.user!.soba!).subscribe({
            next: (soba: Soba) => {
              console.log('Dohvaćena soba:', soba);
              this.roomNumber = soba.roomNumber;
            },
            error: (err) => {
              console.error("Greška pri dohvatanju sobe", err);
              this.roomNumber = 'Nije useljen';
            }
          });
        } else {
          this.roomNumber = 'Nije useljen';
        }
      },
      error: (err) => console.error('Greška pri dohvatanju korisnika', err)
    });
  }
}


  toggleDropdown() { this.dropdownOpen = !this.dropdownOpen; }

  logout(event?: Event) {
    if(event) event.stopPropagation();
    localStorage.removeItem('token');
    this.router.navigate(['/domovi/login']);
  }

  goToProfile(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/domovi/profile']);
  }

  goHome() {
    this.router.navigate(['/domovi/home']);
  }
}