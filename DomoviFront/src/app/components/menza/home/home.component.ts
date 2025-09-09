import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent2 implements OnInit {
  role: string = '';
  userId: string = '';
  dropdownOpen: boolean = false;

  constructor(private router: Router) {}

  ngOnInit(): void {
    // 1. Uzmi token iz localStorage
    const token = localStorage.getItem('token');
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1])); 
      this.role = payload.role;
      console.log(payload)
      this.userId = payload.userId;
    }
  }

  toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  logout(event?: Event) {
    if (event) event.stopPropagation(); // sprečava zatvaranje dropdown-a
    localStorage.removeItem('token');
    localStorage.removeItem('jwt');
    this.router.navigate(['/menza/login']);
  }

  goToProfile(event: Event) {
    event.stopPropagation();
    this.router.navigate(['/menza/profile']); // vodi na profile komponentu
  }
}
