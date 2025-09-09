import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { JeloService } from '../../../../services/jela.service';
import { AuthService } from '../../../../services/auth.service';
import { Jelo } from '../../../../model/Jelo';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router'; 
@Component({
  selector: 'app-prikaz-jela',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule], 
  templateUrl: './prikaz-jela.html',
  styleUrls: ['./prikaz-jela.css']
})
export class PrikazJela implements OnInit {
  jela: Jelo[] = [];
  loading = true;
  tipPretrage: string = '';
  isAdmin = false; // vrednost iz input polja

  constructor(private jeloService: JeloService,private authService:AuthService) {}

  ngOnInit(): void {
    this.loadAllJela();
    this.isAdmin = this.authService.isAdmin();
  }

  loadAllJela(): void {
    this.loading = true;
    this.jeloService.getJela().subscribe({
      next: (data) => { this.jela = data; this.loading = false; },
      error: (err) => { console.error(err); this.loading = false; }
    });
  }

 searchByTip(): void {
  this.loading = true;
  if (!this.tipPretrage) {
    // ako je prazno, učitaj sve
    this.jeloService.getJela().subscribe({
      next: (data) => { this.jela = data; this.loading = false; },
      error: (err) => { console.error(err); this.loading = false; }
    });
    return;
  }

  this.jeloService.getJelaByTip(this.tipPretrage).subscribe({
    next: (data) => { this.jela = data; this.loading = false; },
    error: (err) => { console.error(err); this.loading = false; }
  });
}
}