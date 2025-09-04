import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { JelovnikService } from './../../services/jelovnik.service';
import { JelovnikPrikaz } from './../../model/JelovnikPrikaz';

@Component({
  selector: 'app-prikaz-jelovnik',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './prikaz-jelovnik.html',
  styleUrls: ['./prikaz-jelovnik.css']
})
export class PrikazJelovnik implements OnInit {
  jelovnici: JelovnikPrikaz[] = [];
  loading = true;

  constructor(private jelovnikService: JelovnikService) {}

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
  }
}
