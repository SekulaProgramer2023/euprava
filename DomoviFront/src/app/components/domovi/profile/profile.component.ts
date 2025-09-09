import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
import { UserService } from '../../../services/user.service';
import { User } from '../../../model/user.model';

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

  constructor(private userService: UserService, private router: Router) {}

  ngOnInit(): void {
  const userId = this.userService.getUserIdFromToken();
  if (userId) {
    this.userService.getUserById(userId).subscribe({
      next: (res) => this.user = res,
      error: (err) => console.error('Greška pri dohvatanju korisnika', err)
    });
  }

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

  
  goHome() {
    this.router.navigate(['/domovi/home']);
  }

}
