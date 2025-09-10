import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router'; // <-- obavezno za routerLink
import { JelovnikService } from './../../../../services/jelovnik.service';
import { JelovnikPrikaz } from './../../../../model/JelovnikPrikaz';
import { AuthService } from './../../../../services/auth.service';

@Component({
  selector: 'app-prikaz-jelovnik',
  standalone: true,
  imports: [CommonModule, RouterModule], // <-- dodaj RouterModule ovde
  templateUrl: './prikaz-jelovnik.html',
  styleUrls: ['./prikaz-jelovnik.css']
})
export class PrikazJelovnik implements OnInit {
  jelovnici: JelovnikPrikaz[] = [];
  loading = true;
  isAdmin = false; 

  constructor(
    private jelovnikService: JelovnikService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    this.jelovnikService.getJelovnici().subscribe({
      next: (data) => {
        this.jelovnici = data;
        this.loading = false;
      },
      error: (err) => {
        console.error(err);
        this.loading = false;
      }
    });

    this.isAdmin = this.authService.isAdmin();
    console.log(this.isAdmin)
  }
}
