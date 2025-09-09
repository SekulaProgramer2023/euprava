import { Component, OnInit } from '@angular/core';
import { RoomService, Soba } from '../../../services/room.service';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  sobe: any[] = [];
  role: string = '';
  userId: string = '';
  dropdownOpen: boolean = false;

  constructor(private roomService: RoomService, private router: Router) {}

  ngOnInit(): void {
    // 1. Uzmi token iz localStorage
    const token = localStorage.getItem('token');
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1])); 
      this.role = payload.role;
      console.log(payload)
      this.userId = payload.userId;
    }

    // 2. Učitaj sobe
    this.roomService.getSobe().subscribe({
      next: (data) => this.sobe = data,
      error: (err) => console.error(err)
    });
  }

  useliStudenta(roomId: string): void {
          console.log(roomId, this.userId)
    this.roomService.useliStudenta(roomId, this.userId).subscribe({
      next: (res) => {
        console.log('Uspesno useljen', res);
        // Osvježi listu soba
        this.roomService.getSobe().subscribe(data => this.sobe = data);
      },
      error: (err) => console.error('Greška pri useljavanju', err)

    });
  }

  toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  logout(event?: Event) {
    if(event) event.stopPropagation(); // sprečava zatvaranje dropdown-a
    localStorage.removeItem('token');
    this.router.navigate(['/domovi/login']);
  }

  goToProfile(event: Event) {
  event.stopPropagation();
  this.router.navigate(['/domovi/profile']); // vodi na profile komponentu
}
}

